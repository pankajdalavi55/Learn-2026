# Chapter 6 — Build an MCP Server from Scratch

In this chapter we build a real MCP server step by step. We call it **`mcp-task-server`**. It manages tasks (like a to-do list) and saves them in a small SQLite database.

We write it in **TypeScript** with **Node.js** and the **official MCP SDK** (`@modelcontextprotocol/sdk`). **[VERIFY]** the package name and version against the live docs.

**Our plan:** start tiny, then add more. We also keep the "task logic" separate from the "MCP logic". This makes testing easy.

> **Versions used** (check they are still current):
> - Node.js 20 or newer
> - TypeScript 5.4+
> - `@modelcontextprotocol/sdk` 1.x **[VERIFY]**
> - `zod` for input checking
> - `better-sqlite3` for the database

---

## 6.1 What the server will do

It will have these **tools**:

| Tool | What it does |
|------|--------------|
| `create_task` | Make a new task |
| `list_tasks` | Show all tasks |
| `get_task` | Get one task by its ID |
| `update_task` | Change a task |
| `delete_task` | Remove a task |
| `search_tasks` | Find tasks by text |
| `task_stats` | Show counts (open, done, etc.) |

It will also have:
- One **resource**: `task://stats/summary` (live stats as JSON)
- One **prompt**: `triage_task` (asks the model to rank a task)

---

## 6.2 Project folder structure

Here is the final layout. We build it piece by piece.

```
mcp-task-server/
├── package.json
├── tsconfig.json
├── .env.example
├── .gitignore
├── src/
│   ├── config.ts          # reads and checks settings
│   ├── logger.ts          # safe logging (to stderr)
│   ├── db.ts              # database setup
│   ├── errors.ts          # custom error types
│   ├── tasks/
│   │   ├── schemas.ts     # Zod input shapes
│   │   └── service.ts     # task logic (no MCP here!)
│   ├── mcp/
│   │   ├── server.ts      # builds the MCP server, registers tools
│   │   ├── tools.ts       # tool handlers
│   │   ├── resources.ts   # resource handlers
│   │   └── prompts.ts     # prompt handlers
│   └── index.ts           # entry point (STDIO)
└── tests/
    └── service.test.ts    # tests for task logic
```

**Why split `tasks/` and `mcp/`?**
The `tasks/` folder knows nothing about MCP. It just creates, reads, and deletes tasks. The `mcp/` folder is the glue that connects MCP to the task logic. This means we can test the task logic alone, with no protocol at all.

---

## 6.3 Step 1 — Project setup

### `package.json`
```json
{
  "name": "mcp-task-server",
  "version": "1.0.0",
  "description": "An MCP server that manages tasks.",
  "type": "module",
  "bin": {
    "mcp-task-server": "./build/index.js"
  },
  "scripts": {
    "build": "tsc",
    "start": "node build/index.js",
    "dev": "tsc --watch",
    "test": "vitest run"
  },
  "dependencies": {
    "@modelcontextprotocol/sdk": "^1.0.0",
    "better-sqlite3": "^11.0.0",
    "zod": "^3.23.0"
  },
  "devDependencies": {
    "@types/better-sqlite3": "^7.6.0",
    "@types/node": "^20.0.0",
    "typescript": "^5.4.0",
    "vitest": "^2.0.0"
  }
}
```

**Notes:**
- `"type": "module"` → we use modern ES modules (`import`/`export`).
- `bin` lets other tools run the server by name later.
- **[VERIFY]** the SDK version `^1.0.0`. Match it to the current release.

### `tsconfig.json`
```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "Node16",
    "moduleResolution": "Node16",
    "outDir": "./build",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "declaration": false,
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "build", "tests"]
}
```

### `.gitignore`
```
node_modules/
build/
*.db
.env
```

### Install everything
```bash
npm install
```

---

## 6.4 Step 2 — Config and logging

### `src/config.ts`
This reads settings from environment variables and checks them at startup. If something is wrong, we stop early with a clear message.

```ts
// src/config.ts
// Reads settings from the environment and checks them once at startup.

export interface AppConfig {
  databasePath: string;
  logLevel: "debug" | "info" | "warn" | "error";
  serverName: string;
  serverVersion: string;
}

function readLogLevel(value: string | undefined): AppConfig["logLevel"] {
  const allowed = ["debug", "info", "warn", "error"] as const;
  if (value && (allowed as readonly string[]).includes(value)) {
    return value as AppConfig["logLevel"];
  }
  return "info";
}

export function loadConfig(): AppConfig {
  const config: AppConfig = {
    databasePath: process.env.DATABASE_PATH ?? "./tasks.db",
    logLevel: readLogLevel(process.env.LOG_LEVEL),
    serverName: "mcp-task-server",
    serverVersion: "1.0.0",
  };

  // Simple validation. Fail fast if something is clearly wrong.
  if (config.databasePath.trim() === "") {
    throw new Error("DATABASE_PATH must not be empty.");
  }

  return config;
}
```

### `src/logger.ts`
**Very important for STDIO:** we must **never** write logs to stdout. Stdout is only for MCP messages. So we log to **stderr**.

```ts
// src/logger.ts
// Safe logger. Writes ONLY to stderr, never stdout.
// Stdout is reserved for MCP JSON-RPC messages.

type Level = "debug" | "info" | "warn" | "error";

const order: Record<Level, number> = { debug: 0, info: 1, warn: 2, error: 3 };

export function createLogger(minLevel: Level) {
  function log(level: Level, message: string, extra?: unknown) {
    if (order[level] < order[minLevel]) return;

    const line = {
      time: new Date().toISOString(),
      level,
      message,
      ...(extra !== undefined ? { extra } : {}),
    };

    // console.error writes to stderr. This is safe.
    console.error(JSON.stringify(line));
  }

  return {
    debug: (m: string, e?: unknown) => log("debug", m, e),
    info: (m: string, e?: unknown) => log("info", m, e),
    warn: (m: string, e?: unknown) => log("warn", m, e),
    error: (m: string, e?: unknown) => log("error", m, e),
  };
}

export type Logger = ReturnType<typeof createLogger>;
```

---

## 6.5 Step 3 — Errors

We make our own error types. This helps us tell the difference between "bad input", "not found", and "something broke".

### `src/errors.ts`
```ts
// src/errors.ts
// Our own error types. They make error handling clear and testable.

export class ValidationError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "ValidationError";
  }
}

export class NotFoundError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "NotFoundError";
  }
}

export class DatabaseError extends Error {
  constructor(message: string) {
    super(message);
    this.name = "DatabaseError";
  }
}
```

---

## 6.6 Step 4 — The database

### `src/db.ts`
This sets up SQLite and makes the `tasks` table if it does not exist.

```ts
// src/db.ts
// Opens the SQLite database and creates the tasks table.

import Database from "better-sqlite3";
import type { Logger } from "./logger.js";

export function openDatabase(path: string, logger: Logger): Database.Database {
  const db = new Database(path);

  // WAL mode is safer and faster for many small writes.
  db.pragma("journal_mode = WAL");

  db.exec(`
    CREATE TABLE IF NOT EXISTS tasks (
      id          INTEGER PRIMARY KEY AUTOINCREMENT,
      title       TEXT NOT NULL,
      description TEXT NOT NULL DEFAULT '',
      status      TEXT NOT NULL DEFAULT 'open',
      priority    TEXT NOT NULL DEFAULT 'medium',
      created_at  TEXT NOT NULL,
      updated_at  TEXT NOT NULL
    );
  `);

  logger.info("Database ready", { path });
  return db;
}
```

---

## 6.7 Step 5 — Input shapes (Zod)

### `src/tasks/schemas.ts`
These shapes describe what each tool accepts. Zod checks the data for us. The SDK can also turn these into JSON Schema for the model. **[SDK]**

```ts
// src/tasks/schemas.ts
// Zod shapes for tool inputs. We validate every input with these.

import { z } from "zod";

export const Priority = z.enum(["low", "medium", "high"]);
export const Status = z.enum(["open", "in_progress", "done"]);

export const CreateTaskInput = z.object({
  title: z.string().min(1, "title is required").max(200),
  description: z.string().max(2000).optional(),
  priority: Priority.optional(),
});

export const UpdateTaskInput = z.object({
  id: z.number().int().positive(),
  title: z.string().min(1).max(200).optional(),
  description: z.string().max(2000).optional(),
  status: Status.optional(),
  priority: Priority.optional(),
});

export const GetTaskInput = z.object({
  id: z.number().int().positive(),
});

export const DeleteTaskInput = z.object({
  id: z.number().int().positive(),
});

export const ListTasksInput = z.object({
  status: Status.optional(),
  limit: z.number().int().min(1).max(100).optional(),
});

export const SearchTasksInput = z.object({
  query: z.string().min(1).max(200),
  limit: z.number().int().min(1).max(100).optional(),
});

// Types we can use in the service.
export type CreateTask = z.infer<typeof CreateTaskInput>;
export type UpdateTask = z.infer<typeof UpdateTaskInput>;
export type ListTasks = z.infer<typeof ListTasksInput>;
export type SearchTasks = z.infer<typeof SearchTasksInput>;
```

---

## 6.8 Step 6 — Task logic (no MCP here)

This is the heart. It only knows about tasks and the database. It does **not** know about MCP. This is why we can test it alone.

### `src/tasks/service.ts`
```ts
// src/tasks/service.ts
// All task logic lives here. No MCP code. Easy to test.

import type Database from "better-sqlite3";
import { NotFoundError } from "../errors.js";
import type {
  CreateTask,
  UpdateTask,
  ListTasks,
  SearchTasks,
} from "./schemas.js";

export interface Task {
  id: number;
  title: string;
  description: string;
  status: "open" | "in_progress" | "done";
  priority: "low" | "medium" | "high";
  created_at: string;
  updated_at: string;
}

export interface TaskStats {
  total: number;
  open: number;
  in_progress: number;
  done: number;
}

export class TaskService {
  constructor(private db: Database.Database) {}

  create(input: CreateTask): Task {
    const now = new Date().toISOString();
    const stmt = this.db.prepare(
      `INSERT INTO tasks (title, description, status, priority, created_at, updated_at)
       VALUES (@title, @description, @status, @priority, @created_at, @updated_at)`,
    );
    const result = stmt.run({
      title: input.title,
      description: input.description ?? "",
      status: "open",
      priority: input.priority ?? "medium",
      created_at: now,
      updated_at: now,
    });
    return this.getById(Number(result.lastInsertRowid));
  }

  getById(id: number): Task {
    const row = this.db
      .prepare(`SELECT * FROM tasks WHERE id = ?`)
      .get(id) as Task | undefined;
    if (!row) {
      throw new NotFoundError(`Task ${id} not found`);
    }
    return row;
  }

  list(input: ListTasks): Task[] {
    const limit = input.limit ?? 50;
    if (input.status) {
      return this.db
        .prepare(
          `SELECT * FROM tasks WHERE status = ? ORDER BY id DESC LIMIT ?`,
        )
        .all(input.status, limit) as Task[];
    }
    return this.db
      .prepare(`SELECT * FROM tasks ORDER BY id DESC LIMIT ?`)
      .all(limit) as Task[];
  }

  update(input: UpdateTask): Task {
    // Make sure it exists first. Throws NotFoundError if missing.
    const existing = this.getById(input.id);

    const updated = {
      title: input.title ?? existing.title,
      description: input.description ?? existing.description,
      status: input.status ?? existing.status,
      priority: input.priority ?? existing.priority,
      updated_at: new Date().toISOString(),
    };

    this.db
      .prepare(
        `UPDATE tasks
         SET title=@title, description=@description, status=@status,
             priority=@priority, updated_at=@updated_at
         WHERE id=@id`,
      )
      .run({ ...updated, id: input.id });

    return this.getById(input.id);
  }

  delete(id: number): void {
    const result = this.db.prepare(`DELETE FROM tasks WHERE id = ?`).run(id);
    if (result.changes === 0) {
      throw new NotFoundError(`Task ${id} not found`);
    }
  }

  search(input: SearchTasks): Task[] {
    const limit = input.limit ?? 50;
    // Use a parameter to avoid SQL injection. The % is added safely.
    const like = `%${input.query}%`;
    return this.db
      .prepare(
        `SELECT * FROM tasks
         WHERE title LIKE ? OR description LIKE ?
         ORDER BY id DESC LIMIT ?`,
      )
      .all(like, like, limit) as Task[];
  }

  stats(): TaskStats {
    const rows = this.db
      .prepare(`SELECT status, COUNT(*) as count FROM tasks GROUP BY status`)
      .all() as { status: string; count: number }[];

    const stats: TaskStats = { total: 0, open: 0, in_progress: 0, done: 0 };
    for (const row of rows) {
      stats.total += row.count;
      if (row.status === "open") stats.open = row.count;
      if (row.status === "in_progress") stats.in_progress = row.count;
      if (row.status === "done") stats.done = row.count;
    }
    return stats;
  }
}
```

**Notice:** every SQL query uses `?` or named parameters. We never build SQL by joining strings with user input. This stops SQL injection.

---

## 6.9 Step 7 — The MCP tools

Now we connect MCP to the task logic. Each tool:
1. Gets arguments from the model.
2. Checks them with Zod.
3. Calls the service.
4. Returns a result.

> **[SDK/VERIFY]** The exact method to register a tool (`server.registerTool(...)` vs `server.tool(...)`) depends on the SDK version. Below uses `registerTool` with a Zod shape. If your SDK uses a different name, change it. The *idea* stays the same.

### `src/mcp/tools.ts`
```ts
// src/mcp/tools.ts
// Connects MCP tool calls to the TaskService.
// Checks input, calls the service, returns MCP content.

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { z } from "zod";
import type { TaskService, Task } from "../tasks/service.js";
import { NotFoundError, ValidationError } from "../errors.js";
import type { Logger } from "../logger.js";
import {
  CreateTaskInput,
  UpdateTaskInput,
  GetTaskInput,
  DeleteTaskInput,
  ListTasksInput,
  SearchTasksInput,
} from "../tasks/schemas.js";

// Turns any value into an MCP text result.
function textResult(text: string, isError = false) {
  return {
    content: [{ type: "text" as const, text }],
    isError,
  };
}

// Shows a task as readable text.
function formatTask(t: Task): string {
  return `#${t.id} [${t.status}] (${t.priority}) ${t.title}`;
}

// Wraps a handler so business errors become nice MCP error results,
// not crashes.
function safe(
  logger: Logger,
  handler: (args: unknown) => Promise<ReturnType<typeof textResult>>,
) {
  return async (args: unknown) => {
    try {
      return await handler(args);
    } catch (err) {
      if (err instanceof NotFoundError || err instanceof ValidationError) {
        // Expected problem. Tell the model what went wrong.
        return textResult(err.message, true);
      }
      // Unexpected problem. Log details, but do not leak them.
      logger.error("Tool failed", { error: String(err) });
      return textResult("An internal error occurred.", true);
    }
  };
}

export function registerTools(
  server: McpServer,
  service: TaskService,
  logger: Logger,
) {
  server.registerTool(
    "create_task",
    {
      description:
        "Create a new task. Use this when the user wants to add a to-do item.",
      inputSchema: CreateTaskInput.shape,
    },
    safe(logger, async (args) => {
      const input = CreateTaskInput.parse(args);
      const task = service.create(input);
      return textResult(`Created task: ${formatTask(task)}`);
    }),
  );

  server.registerTool(
    "list_tasks",
    {
      description: "List tasks. You can filter by status and set a limit.",
      inputSchema: ListTasksInput.shape,
    },
    safe(logger, async (args) => {
      const input = ListTasksInput.parse(args);
      const tasks = service.list(input);
      if (tasks.length === 0) return textResult("No tasks found.");
      return textResult(tasks.map(formatTask).join("\n"));
    }),
  );

  server.registerTool(
    "get_task",
    {
      description: "Get one task by its numeric ID.",
      inputSchema: GetTaskInput.shape,
    },
    safe(logger, async (args) => {
      const input = GetTaskInput.parse(args);
      const task = service.getById(input.id);
      const full = `${formatTask(task)}\n${task.description}`;
      return textResult(full);
    }),
  );

  server.registerTool(
    "update_task",
    {
      description:
        "Update a task. Give the ID and the fields you want to change.",
      inputSchema: UpdateTaskInput.shape,
    },
    safe(logger, async (args) => {
      const input = UpdateTaskInput.parse(args);
      const task = service.update(input);
      return textResult(`Updated task: ${formatTask(task)}`);
    }),
  );

  server.registerTool(
    "delete_task",
    {
      description: "Delete a task by its ID. This cannot be undone.",
      inputSchema: DeleteTaskInput.shape,
    },
    safe(logger, async (args) => {
      const input = DeleteTaskInput.parse(args);
      service.delete(input.id);
      return textResult(`Deleted task #${input.id}.`);
    }),
  );

  server.registerTool(
    "search_tasks",
    {
      description: "Search tasks by text in the title or description.",
      inputSchema: SearchTasksInput.shape,
    },
    safe(logger, async (args) => {
      const input = SearchTasksInput.parse(args);
      const tasks = service.search(input);
      if (tasks.length === 0) return textResult("No matching tasks.");
      return textResult(tasks.map(formatTask).join("\n"));
    }),
  );

  server.registerTool(
    "task_stats",
    {
      description: "Show counts of tasks by status.",
      inputSchema: z.object({}).shape,
    },
    safe(logger, async () => {
      const s = service.stats();
      return textResult(
        `Total: ${s.total}, Open: ${s.open}, In progress: ${s.in_progress}, Done: ${s.done}`,
      );
    }),
  );

  logger.info("Tools registered");
}
```

**Key points:**
- We call `.parse(args)` to check input again inside the handler. The schema we gave the model is a *hint*. We still validate. **Never trust the model's arguments.**
- Business errors (`NotFoundError`) become `isError: true` results. The model can read them and react.
- Unexpected errors get logged but **not** shown to the model in detail. This avoids leaking secrets.

---

## 6.10 Step 8 — The resource

### `src/mcp/resources.ts`
```ts
// src/mcp/resources.ts
// One resource: live task stats as JSON.

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import type { TaskService } from "../tasks/service.js";
import type { Logger } from "../logger.js";

export function registerResources(
  server: McpServer,
  service: TaskService,
  logger: Logger,
) {
  server.registerResource(
    "task-stats",
    "task://stats/summary",
    {
      title: "Task statistics",
      description: "Live counts of tasks by status, as JSON.",
      mimeType: "application/json",
    },
    async (uri) => {
      const stats = service.stats();
      return {
        contents: [
          {
            uri: uri.href,
            mimeType: "application/json",
            text: JSON.stringify(stats, null, 2),
          },
        ],
      };
    },
  );

  logger.info("Resources registered");
}
```

> **[SDK/VERIFY]** `registerResource` arguments and the URI object shape depend on the SDK version. Check the current API.

---

## 6.11 Step 9 — The prompt

### `src/mcp/prompts.ts`
```ts
// src/mcp/prompts.ts
// One prompt: triage_task. It builds a message asking the model
// to rank a task by urgency.

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { z } from "zod";
import type { TaskService } from "../tasks/service.js";
import type { Logger } from "../logger.js";

export function registerPrompts(
  server: McpServer,
  service: TaskService,
  logger: Logger,
) {
  server.registerPrompt(
    "triage_task",
    {
      title: "Triage a task",
      description: "Ask the model to rank one task by urgency.",
      argsSchema: { taskId: z.string() },
    },
    async (args) => {
      const id = Number(args.taskId);
      // We still validate. The prompt argument is user input.
      const task = service.getById(id);

      return {
        messages: [
          {
            role: "user",
            content: {
              type: "text",
              text:
                `Please triage this task. Suggest a priority (low/medium/high) ` +
                `and a short reason.\n\n` +
                `Title: ${task.title}\n` +
                `Description: ${task.description || "(none)"}\n` +
                `Current status: ${task.status}\n` +
                `Current priority: ${task.priority}`,
            },
          },
        ],
      };
    },
  );

  logger.info("Prompts registered");
}
```

> **[SDK/VERIFY]** `registerPrompt` and `argsSchema` names depend on the SDK version. Check the current API.

---

## 6.12 Step 10 — Build the server

This file builds the MCP server object and registers everything. It does **not** choose a transport yet. That keeps it reusable for STDIO or HTTP later.

### `src/mcp/server.ts`
```ts
// src/mcp/server.ts
// Builds the MCP server and registers tools, resources, prompts.
// No transport here. The entry point picks the transport.

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import type { TaskService } from "../tasks/service.js";
import type { Logger } from "../logger.js";
import type { AppConfig } from "../config.js";
import { registerTools } from "./tools.js";
import { registerResources } from "./resources.js";
import { registerPrompts } from "./prompts.js";

export function buildServer(
  config: AppConfig,
  service: TaskService,
  logger: Logger,
): McpServer {
  const server = new McpServer({
    name: config.serverName,
    version: config.serverVersion,
  });

  registerTools(server, service, logger);
  registerResources(server, service, logger);
  registerPrompts(server, service, logger);

  return server;
}
```

---

## 6.13 Step 11 — The entry point (STDIO)

This is the file that runs. It wires everything together and uses the STDIO transport.

### `src/index.ts`
```ts
// src/index.ts
// The program that runs. Uses STDIO transport (local use).

import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { loadConfig } from "./config.js";
import { createLogger } from "./logger.js";
import { openDatabase } from "./db.js";
import { TaskService } from "./tasks/service.js";
import { buildServer } from "./mcp/server.js";

async function main() {
  const config = loadConfig();
  const logger = createLogger(config.logLevel);

  logger.info("Starting mcp-task-server", { version: config.serverVersion });

  const db = openDatabase(config.databasePath, logger);
  const service = new TaskService(db);
  const server = buildServer(config, service, logger);

  const transport = new StdioServerTransport();

  // Clean shutdown. Close the DB and exit.
  async function shutdown() {
    logger.info("Shutting down");
    try {
      db.close();
    } catch (err) {
      logger.error("Error closing database", { error: String(err) });
    }
    process.exit(0);
  }

  process.on("SIGINT", shutdown);
  process.on("SIGTERM", shutdown);
  // When the host closes stdin, we should exit too.
  process.stdin.on("end", shutdown);

  await server.connect(transport);
  logger.info("Server connected over STDIO");
}

main().catch((err) => {
  // Last line of defense. Log to stderr and exit non-zero.
  console.error("Fatal error:", err);
  process.exit(1);
});
```

> **[SDK/VERIFY]** Import paths like `.../server/stdio.js` and `.../server/mcp.js` depend on the SDK version. Check the current package exports.

---

## 6.14 Step 12 — Tests for the logic

Because the task logic has no MCP code, we can test it with a real (temporary) database and no protocol at all.

### `tests/service.test.ts`
```ts
// tests/service.test.ts
// Tests the TaskService with an in-memory SQLite database.

import { describe, it, expect, beforeEach } from "vitest";
import Database from "better-sqlite3";
import { TaskService } from "../src/tasks/service.js";
import { NotFoundError } from "../src/errors.js";

function makeDb(): Database.Database {
  const db = new Database(":memory:");
  db.exec(`
    CREATE TABLE tasks (
      id          INTEGER PRIMARY KEY AUTOINCREMENT,
      title       TEXT NOT NULL,
      description TEXT NOT NULL DEFAULT '',
      status      TEXT NOT NULL DEFAULT 'open',
      priority    TEXT NOT NULL DEFAULT 'medium',
      created_at  TEXT NOT NULL,
      updated_at  TEXT NOT NULL
    );
  `);
  return db;
}

describe("TaskService", () => {
  let service: TaskService;

  beforeEach(() => {
    service = new TaskService(makeDb());
  });

  it("creates and gets a task", () => {
    const task = service.create({ title: "Test task" });
    expect(task.id).toBeGreaterThan(0);
    expect(task.title).toBe("Test task");
    expect(task.status).toBe("open");

    const same = service.getById(task.id);
    expect(same.title).toBe("Test task");
  });

  it("throws when a task is missing", () => {
    expect(() => service.getById(999)).toThrow(NotFoundError);
  });

  it("updates a task", () => {
    const task = service.create({ title: "Old" });
    const updated = service.update({ id: task.id, title: "New", status: "done" });
    expect(updated.title).toBe("New");
    expect(updated.status).toBe("done");
  });

  it("deletes a task", () => {
    const task = service.create({ title: "Remove me" });
    service.delete(task.id);
    expect(() => service.getById(task.id)).toThrow(NotFoundError);
  });

  it("searches tasks", () => {
    service.create({ title: "Buy milk" });
    service.create({ title: "Buy bread" });
    service.create({ title: "Call mom" });
    const results = service.search({ query: "Buy" });
    expect(results.length).toBe(2);
  });

  it("counts stats", () => {
    service.create({ title: "A" });
    const t = service.create({ title: "B" });
    service.update({ id: t.id, status: "done" });
    const stats = service.stats();
    expect(stats.total).toBe(2);
    expect(stats.open).toBe(1);
    expect(stats.done).toBe(1);
  });
});
```

Run the tests:
```bash
npm test
```

---

## 6.15 Step 13 — Build and run

```bash
# Build TypeScript to JavaScript
npm run build

# Run the server
npm start
```

The server now waits for MCP messages on stdin. If you type nothing, nothing happens — that is normal. A real client (or the MCP Inspector) talks to it. We cover that in Chapter 7 and Chapter 8.

### Quick manual check with the Inspector
```bash
npx @modelcontextprotocol/inspector node build/index.js
```
This opens a tool where you can see the tools, call them, read the resource, and try the prompt. **[VERIFY]** the Inspector package name and usage.

---

## 6.16 What we built and why

| Choice | Reason |
|--------|--------|
| Split `tasks/` from `mcp/` | Test logic without the protocol |
| Zod + `.parse()` inside handlers | Never trust model input |
| Parameterized SQL | Stop SQL injection |
| Logs to stderr only | Keep stdout clean for STDIO |
| Custom error types | Clear, testable error handling |
| `isError: true` for business failures | Let the model see and react |
| Transport chosen in `index.ts` | Easy to add HTTP later |

### Common mistakes to avoid
- Using `console.log` → it breaks STDIO. Use `console.error` or the MCP log channel.
- Skipping `.parse()` because "the schema already checks it" → the schema is only a hint to the model.
- Building SQL with string joins → injection risk.
- Returning raw internal errors to the model → can leak secrets.
- Putting auth logic only in the tool description → descriptions are not security.

---