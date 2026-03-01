# JUnit 5 & Mockito — Complete Testing Guide

> A comprehensive guide covering unit testing, integration testing, and mocking in Spring Boot applications.  
> **Part 7 of the Spring Boot & Microservices Series**  
> **Prerequisites:** Spring Boot basics, Java fundamentals

---

## Table of Contents

1. [JUnit 5 Fundamentals](#1-junit-5-fundamentals)
2. [Mockito Core Concepts](#2-mockito-core-concepts)
3. [Spring Boot Test Integration](#3-spring-boot-test-integration)
4. [MockMvc & REST API Testing](#4-mockmvc--rest-api-testing)
5. [Database Testing with Testcontainers](#5-database-testing-with-testcontainers)
6. [Advanced Mocking Patterns](#6-advanced-mocking-patterns)
7. [Testing Best Practices & Patterns](#7-testing-best-practices--patterns)
8. [Interview Questions](#8-interview-questions)
9. [Complete CRUD Testing Example](#9-complete-crud-testing-example)

---

## 1. JUnit 5 Fundamentals

### 1.1 JUnit 5 Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         JUnit 5 Architecture                                     │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                         JUnit Platform                                   │   │
│  │  • Foundation for launching testing frameworks                          │   │
│  │  • TestEngine API for building test frameworks                          │   │
│  │  • Console Launcher, IDE & build tool support                           │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                             │
│              ┌─────────────────────┼─────────────────────┐                      │
│              │                     │                     │                      │
│              ▼                     ▼                     ▼                      │
│  ┌─────────────────────┐  ┌─────────────────┐  ┌─────────────────────┐         │
│  │   JUnit Jupiter     │  │  JUnit Vintage  │  │  Third-party        │         │
│  │  (JUnit 5 tests)    │  │ (JUnit 3/4)     │  │  (Spock, etc.)      │         │
│  │                     │  │                 │  │                     │         │
│  │  • New programming  │  │  • Backward     │  │  • Custom engines   │         │
│  │    model            │  │    compatibility│  │                     │         │
│  │  • Extension model  │  │                 │  │                     │         │
│  └─────────────────────┘  └─────────────────┘  └─────────────────────┘         │
│                                                                                  │
│  DEPENDENCY:                                                                     │
│  <dependency>                                                                    │
│      <groupId>org.junit.jupiter</groupId>                                       │
│      <artifactId>junit-jupiter</artifactId>                                     │
│      <scope>test</scope>                                                        │
│  </dependency>                                                                   │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Basic Annotations

```java
import org.junit.jupiter.api.*;
import static org.junit.jupiter.api.Assertions.*;

class CalculatorTest {

    private Calculator calculator;

    @BeforeAll
    static void initAll() {
        // Runs once before all tests (must be static)
        System.out.println("Starting test suite...");
    }

    @BeforeEach
    void init() {
        // Runs before each test
        calculator = new Calculator();
    }

    @Test
    @DisplayName("Adding two positive numbers")
    void addPositiveNumbers() {
        assertEquals(5, calculator.add(2, 3), "2 + 3 should equal 5");
    }

    @Test
    @DisplayName("Division by zero should throw exception")
    void divisionByZero() {
        ArithmeticException exception = assertThrows(
            ArithmeticException.class,
            () -> calculator.divide(10, 0),
            "Division by zero should throw"
        );
        assertEquals("/ by zero", exception.getMessage());
    }

    @Test
    @Disabled("Temporarily disabled - bug #123")
    void skippedTest() {
        // This test won't run
    }

    @AfterEach
    void tearDown() {
        // Runs after each test
        calculator = null;
    }

    @AfterAll
    static void tearDownAll() {
        // Runs once after all tests (must be static)
        System.out.println("Test suite completed.");
    }
}
```

### 1.3 Assertions

```java
import static org.junit.jupiter.api.Assertions.*;

class AssertionsDemo {

    @Test
    void basicAssertions() {
        // Equality
        assertEquals(4, calculator.add(2, 2));
        assertEquals("JUnit", "JUnit");
        
        // Not equal
        assertNotEquals(5, calculator.add(2, 2));
        
        // Boolean
        assertTrue(calculator.isPositive(5));
        assertFalse(calculator.isNegative(5));
        
        // Null checks
        assertNull(getNullValue());
        assertNotNull(calculator);
        
        // Same reference
        Object obj = new Object();
        assertSame(obj, obj);
        assertNotSame(new Object(), new Object());
    }

    @Test
    void arrayAssertions() {
        int[] expected = {1, 2, 3};
        int[] actual = {1, 2, 3};
        
        assertArrayEquals(expected, actual);
    }

    @Test
    void iterableAssertions() {
        List<String> expected = List.of("a", "b", "c");
        List<String> actual = List.of("a", "b", "c");
        
        assertIterableEquals(expected, actual);
    }

    @Test
    void groupedAssertions() {
        User user = new User("John", "Doe", 30);
        
        // All assertions run, even if some fail
        assertAll("User properties",
            () -> assertEquals("John", user.getFirstName()),
            () -> assertEquals("Doe", user.getLastName()),
            () -> assertTrue(user.getAge() > 18)
        );
    }

    @Test
    void exceptionAssertion() {
        // Assert exception type
        IllegalArgumentException exception = assertThrows(
            IllegalArgumentException.class,
            () -> validateAge(-1)
        );
        
        // Assert exception message
        assertTrue(exception.getMessage().contains("negative"));
        
        // Assert no exception
        assertDoesNotThrow(() -> validateAge(25));
    }

    @Test
    void timeoutAssertion() {
        // Assert completes within timeout
        assertTimeout(Duration.ofSeconds(2), () -> {
            Thread.sleep(500);
            return "completed";
        });
        
        // Preemptively abort if timeout exceeded
        assertTimeoutPreemptively(Duration.ofSeconds(1), () -> {
            // This would be killed if it takes > 1 second
            return computeResult();
        });
    }

    @Test
    void conditionalAssertions() {
        // Fail with message
        if (someCondition()) {
            fail("Test should not reach this point");
        }
    }
}
```

### 1.4 Parameterized Tests

```java
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.*;

class ParameterizedTestsDemo {

    // Using @ValueSource
    @ParameterizedTest
    @ValueSource(ints = {1, 2, 3, 4, 5})
    void testPositiveNumbers(int number) {
        assertTrue(number > 0);
    }

    @ParameterizedTest
    @ValueSource(strings = {"hello", "world", "junit"})
    void testNonEmptyStrings(String input) {
        assertFalse(input.isEmpty());
    }

    // Test with null and empty values
    @ParameterizedTest
    @NullSource
    @EmptySource
    @ValueSource(strings = {"  ", "\t", "\n"})
    void testBlankStrings(String input) {
        assertTrue(input == null || input.isBlank());
    }

    // Combined null and empty
    @ParameterizedTest
    @NullAndEmptySource
    void testNullAndEmpty(String input) {
        assertTrue(input == null || input.isEmpty());
    }

    // Using @EnumSource
    @ParameterizedTest
    @EnumSource(OrderStatus.class)
    void testAllOrderStatuses(OrderStatus status) {
        assertNotNull(status);
    }

    @ParameterizedTest
    @EnumSource(value = OrderStatus.class, names = {"PENDING", "PROCESSING"})
    void testActiveStatuses(OrderStatus status) {
        assertTrue(status.isActive());
    }

    @ParameterizedTest
    @EnumSource(value = OrderStatus.class, mode = EnumSource.Mode.EXCLUDE, 
                names = {"CANCELLED", "FAILED"})
    void testNonFailedStatuses(OrderStatus status) {
        assertFalse(status.isFailed());
    }

    // Using @CsvSource
    @ParameterizedTest
    @CsvSource({
        "1, 1, 2",
        "2, 3, 5",
        "10, 20, 30",
        "-5, 5, 0"
    })
    void testAddition(int a, int b, int expected) {
        assertEquals(expected, calculator.add(a, b));
    }

    @ParameterizedTest
    @CsvSource(value = {
        "John, Doe, John Doe",
        "Jane, Smith, Jane Smith",
        "'', Doe, Doe"  // Empty string
    })
    void testFullName(String first, String last, String expected) {
        assertEquals(expected, createFullName(first, last));
    }

    // Using CSV file
    @ParameterizedTest
    @CsvFileSource(resources = "/test-data.csv", numLinesToSkip = 1)
    void testFromCsvFile(String input, int expected) {
        assertEquals(expected, parseInt(input));
    }

    // Using @MethodSource
    @ParameterizedTest
    @MethodSource("provideArgumentsForAdd")
    void testAddWithMethod(int a, int b, int expected) {
        assertEquals(expected, calculator.add(a, b));
    }

    static Stream<Arguments> provideArgumentsForAdd() {
        return Stream.of(
            Arguments.of(1, 1, 2),
            Arguments.of(2, 3, 5),
            Arguments.of(100, 200, 300)
        );
    }

    // Complex objects with @MethodSource
    @ParameterizedTest
    @MethodSource("provideUsers")
    void testUserValidation(User user, boolean expected) {
        assertEquals(expected, userService.isValid(user));
    }

    static Stream<Arguments> provideUsers() {
        return Stream.of(
            Arguments.of(new User("John", "john@email.com"), true),
            Arguments.of(new User("", "test@email.com"), false),
            Arguments.of(new User("Jane", "invalid-email"), false)
        );
    }

    // Using ArgumentsProvider
    @ParameterizedTest
    @ArgumentsSource(CustomArgumentsProvider.class)
    void testWithCustomProvider(String input, int expected) {
        assertEquals(expected, input.length());
    }
}

// Custom ArgumentsProvider
class CustomArgumentsProvider implements ArgumentsProvider {
    @Override
    public Stream<? extends Arguments> provideArguments(ExtensionContext context) {
        return Stream.of(
            Arguments.of("hello", 5),
            Arguments.of("world", 5),
            Arguments.of("JUnit5", 6)
        );
    }
}
```

### 1.5 Test Lifecycle & Execution Order

```java
import org.junit.jupiter.api.*;

@TestInstance(TestInstance.Lifecycle.PER_CLASS) // Non-static @BeforeAll/@AfterAll
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
class OrderedTestsDemo {

    private int counter = 0;

    @BeforeAll
    void init() {
        // With PER_CLASS, this doesn't need to be static
        counter = 0;
    }

    @Test
    @Order(1)
    void firstTest() {
        counter++;
        assertEquals(1, counter);
    }

    @Test
    @Order(2)
    void secondTest() {
        counter++;
        assertEquals(2, counter);
    }

    @Test
    @Order(3)
    void thirdTest() {
        counter++;
        assertEquals(3, counter);
    }
}

// Other ordering strategies
@TestMethodOrder(MethodOrderer.DisplayName.class)     // Alphabetically by display name
@TestMethodOrder(MethodOrderer.MethodName.class)     // Alphabetically by method name
@TestMethodOrder(MethodOrderer.Random.class)         // Random order (default)
```

### 1.6 Nested Tests

```java
@DisplayName("User Service Tests")
class UserServiceTest {

    private UserService userService;

    @BeforeEach
    void setup() {
        userService = new UserService();
    }

    @Nested
    @DisplayName("When creating users")
    class UserCreation {

        @Test
        @DisplayName("should create user with valid data")
        void createValidUser() {
            User user = userService.create("John", "john@email.com");
            assertNotNull(user.getId());
        }

        @Test
        @DisplayName("should throw exception for invalid email")
        void createUserInvalidEmail() {
            assertThrows(ValidationException.class,
                () -> userService.create("John", "invalid"));
        }

        @Nested
        @DisplayName("With existing email")
        class WithExistingEmail {

            @BeforeEach
            void createExistingUser() {
                userService.create("Existing", "existing@email.com");
            }

            @Test
            @DisplayName("should throw duplicate exception")
            void createDuplicateEmail() {
                assertThrows(DuplicateException.class,
                    () -> userService.create("New", "existing@email.com"));
            }
        }
    }

    @Nested
    @DisplayName("When finding users")
    class UserRetrieval {

        private User existingUser;

        @BeforeEach
        void createUser() {
            existingUser = userService.create("Test", "test@email.com");
        }

        @Test
        @DisplayName("should find user by ID")
        void findById() {
            User found = userService.findById(existingUser.getId());
            assertEquals(existingUser.getEmail(), found.getEmail());
        }

        @Test
        @DisplayName("should return empty for non-existent ID")
        void findByInvalidId() {
            Optional<User> found = userService.findById(999L);
            assertTrue(found.isEmpty());
        }
    }
}
```

### 1.7 Conditional Test Execution

```java
import org.junit.jupiter.api.condition.*;

class ConditionalTestsDemo {

    // Operating System conditions
    @Test
    @EnabledOnOs(OS.WINDOWS)
    void onlyOnWindows() {
        // Runs only on Windows
    }

    @Test
    @DisabledOnOs({OS.MAC, OS.LINUX})
    void notOnMacOrLinux() {
        // Runs on any OS except Mac and Linux
    }

    // JRE Version conditions
    @Test
    @EnabledOnJre(JRE.JAVA_21)
    void onlyOnJava21() {
        // Runs only on Java 21
    }

    @Test
    @EnabledForJreRange(min = JRE.JAVA_17, max = JRE.JAVA_21)
    void onJava17To21() {
        // Runs on Java 17 through 21
    }

    // System property conditions
    @Test
    @EnabledIfSystemProperty(named = "env", matches = "test")
    void onlyInTestEnv() {
        // Runs only when -Denv=test
    }

    // Environment variable conditions
    @Test
    @EnabledIfEnvironmentVariable(named = "CI", matches = "true")
    void onlyOnCI() {
        // Runs only on CI server
    }

    // Custom conditions
    @Test
    @EnabledIf("customCondition")
    void conditionalTest() {
        // Runs only if customCondition() returns true
    }

    boolean customCondition() {
        return LocalTime.now().getHour() < 18;
    }

    @Test
    @DisabledIf("isWeekend")
    void notOnWeekends() {
        // Disabled on weekends
    }

    boolean isWeekend() {
        DayOfWeek day = LocalDate.now().getDayOfWeek();
        return day == DayOfWeek.SATURDAY || day == DayOfWeek.SUNDAY;
    }
}
```

### 1.8 Repeated Tests

```java
class RepeatedTestsDemo {

    @RepeatedTest(5)
    void repeatedTest() {
        // Runs 5 times
        assertTrue(Math.random() < 1.0);
    }

    @RepeatedTest(value = 3, name = "{displayName} - repetition {currentRepetition}/{totalRepetitions}")
    @DisplayName("Performance test")
    void repeatedWithCustomName(RepetitionInfo info) {
        System.out.println("Repetition: " + info.getCurrentRepetition() + 
                          " of " + info.getTotalRepetitions());
    }
}
```

---

## 2. Mockito Core Concepts

### 2.1 Why Mocking?

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           Why We Need Mocking                                    │
│                                                                                  │
│  PROBLEM: Unit testing a class with dependencies                                │
│                                                                                  │
│  ┌─────────────────┐                                                            │
│  │  OrderService   │                                                            │
│  │                 │                                                            │
│  │  createOrder()  │──depends on──▶ ┌──────────────┐                           │
│  │                 │                 │ UserClient   │──▶ External API           │
│  │                 │──depends on──▶ │ PaymentGateway──▶ Payment Provider        │
│  │                 │──depends on──▶ │ EmailService │──▶ SMTP Server            │
│  │                 │──depends on──▶ │ Database     │──▶ PostgreSQL             │
│  └─────────────────┘                 └──────────────┘                           │
│                                                                                  │
│  Issues with real dependencies:                                                 │
│  • Slow (network calls, DB queries)                                             │
│  • Unreliable (external services might be down)                                 │
│  • Hard to simulate errors                                                      │
│  • Side effects (sends real emails, charges cards)                              │
│                                                                                  │
│  SOLUTION: Mock the dependencies                                                │
│                                                                                  │
│  ┌─────────────────┐                                                            │
│  │  OrderService   │                                                            │
│  │                 │──uses──▶ ┌──────────────────┐                              │
│  │  createOrder()  │          │   Mock Objects   │                              │
│  │   (real code)   │          │                  │                              │
│  │                 │          │  • Controlled    │                              │
│  │                 │          │  • Fast          │                              │
│  │                 │          │  • Predictable   │                              │
│  └─────────────────┘          │  • No side effects                              │
│                               └──────────────────┘                              │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Mockito Setup

```xml
<!-- pom.xml -->
<dependency>
    <groupId>org.mockito</groupId>
    <artifactId>mockito-core</artifactId>
    <version>5.10.0</version>
    <scope>test</scope>
</dependency>
<dependency>
    <groupId>org.mockito</groupId>
    <artifactId>mockito-junit-jupiter</artifactId>
    <version>5.10.0</version>
    <scope>test</scope>
</dependency>
```

### 2.3 Creating Mocks

```java
import org.mockito.Mock;
import org.mockito.InjectMocks;
import org.mockito.junit.jupiter.MockitoExtension;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class OrderServiceTest {

    // Method 1: @Mock annotation
    @Mock
    private UserRepository userRepository;

    @Mock
    private PaymentGateway paymentGateway;

    // @InjectMocks creates instance and injects mocks
    @InjectMocks
    private OrderService orderService;

    // Method 2: Manual creation
    @Test
    void manualMockCreation() {
        UserRepository mockRepo = mock(UserRepository.class);
        PaymentGateway mockGateway = mock(PaymentGateway.class);
        
        OrderService service = new OrderService(mockRepo, mockGateway);
    }

    // Method 3: Mock with settings
    @Test
    void mockWithSettings() {
        UserRepository mockRepo = mock(UserRepository.class, withSettings()
            .name("mockUserRepo")
            .verboseLogging()
            .defaultAnswer(RETURNS_SMART_NULLS));
    }
}
```

### 2.4 Stubbing Methods

```java
@ExtendWith(MockitoExtension.class)
class StubbingDemo {

    @Mock
    private UserRepository userRepository;

    @Mock
    private PaymentGateway paymentGateway;

    @Test
    void basicStubbing() {
        // Stub method to return specific value
        User mockUser = new User(1L, "John", "john@email.com");
        when(userRepository.findById(1L)).thenReturn(Optional.of(mockUser));

        // Now call returns our stubbed value
        Optional<User> result = userRepository.findById(1L);
        assertEquals("John", result.get().getName());
    }

    @Test
    void stubbingVoidMethods() {
        // For void methods, use doNothing/doThrow
        doNothing().when(userRepository).delete(any());
        
        doThrow(new RuntimeException("DB Error"))
            .when(userRepository).delete(999L);
    }

    @Test
    void stubbingWithThrow() {
        when(userRepository.findById(999L))
            .thenThrow(new UserNotFoundException("Not found"));

        assertThrows(UserNotFoundException.class,
            () -> userRepository.findById(999L));
    }

    @Test
    void stubbingConsecutiveCalls() {
        // Different returns for consecutive calls
        when(paymentGateway.processPayment(any()))
            .thenReturn(Result.PENDING)     // First call
            .thenReturn(Result.APPROVED)    // Second call
            .thenReturn(Result.DECLINED);   // Third+ calls

        assertEquals(Result.PENDING, paymentGateway.processPayment(new Payment()));
        assertEquals(Result.APPROVED, paymentGateway.processPayment(new Payment()));
        assertEquals(Result.DECLINED, paymentGateway.processPayment(new Payment()));
    }

    @Test
    void stubbingWithAnswer() {
        // Dynamic return based on input
        when(userRepository.save(any(User.class)))
            .thenAnswer(invocation -> {
                User user = invocation.getArgument(0);
                user.setId(100L); // Simulate ID generation
                return user;
            });

        User saved = userRepository.save(new User("Jane"));
        assertEquals(100L, saved.getId());
    }

    @Test
    void stubbingWithThenAnswer() {
        // Access method arguments in answer
        when(userRepository.findByEmail(anyString()))
            .thenAnswer(invocation -> {
                String email = invocation.getArgument(0);
                return Optional.of(new User(1L, email.split("@")[0], email));
            });

        Optional<User> user = userRepository.findByEmail("test@example.com");
        assertEquals("test", user.get().getName());
    }
}
```

### 2.5 Argument Matchers

```java
import static org.mockito.ArgumentMatchers.*;

class ArgumentMatchersDemo {

    @Mock
    private UserService userService;

    @Test
    void anyMatchers() {
        // Match any value of type
        when(userService.findById(anyLong())).thenReturn(Optional.empty());
        when(userService.findByName(anyString())).thenReturn(List.of());
        when(userService.save(any(User.class))).thenReturn(new User());
        when(userService.findAll(anyList())).thenReturn(List.of());
        when(userService.findByMap(anyMap())).thenReturn(List.of());
    }

    @Test
    void specificMatchers() {
        // Specific value matching
        when(userService.findById(eq(1L))).thenReturn(Optional.of(user1));
        when(userService.findById(eq(2L))).thenReturn(Optional.of(user2));
    }

    @Test
    void nullMatchers() {
        when(userService.process(isNull())).thenReturn("null input");
        when(userService.process(isNotNull())).thenReturn("non-null input");
    }

    @Test
    void stringMatchers() {
        // String patterns
        when(userService.findByName(startsWith("John"))).thenReturn(List.of(john));
        when(userService.findByName(endsWith("Doe"))).thenReturn(List.of(johnDoe));
        when(userService.findByName(contains("ohn"))).thenReturn(List.of(john));
        when(userService.findByName(matches("\\w+@\\w+\\.\\w+"))).thenReturn(List.of());
    }

    @Test
    void customMatchers() {
        // Custom argument matcher
        when(userService.save(argThat(user -> 
            user.getAge() >= 18 && user.getEmail() != null)))
            .thenReturn(new User());
    }

    @Test
    void combinedMatchers() {
        // IMPORTANT: If using matchers, ALL arguments must use matchers
        // WRONG: userService.method(anyString(), "literal") 
        // CORRECT:
        when(userService.search(anyString(), eq("active"), anyInt()))
            .thenReturn(List.of());
    }
}
```

### 2.6 Verifying Interactions

```java
class VerificationDemo {

    @Mock
    private UserRepository userRepository;
    
    @Mock
    private EmailService emailService;

    @InjectMocks
    private UserService userService;

    @Test
    void basicVerification() {
        User user = new User("John", "john@email.com");
        userService.createUser(user);

        // Verify method was called
        verify(userRepository).save(user);
        verify(emailService).sendWelcomeEmail("john@email.com");
    }

    @Test
    void verifyWithTimes() {
        // Verify exact number of invocations
        verify(userRepository, times(1)).save(any());
        verify(userRepository, times(2)).findById(anyLong());
        
        // Verify never called
        verify(emailService, never()).sendEmail(any());
        verify(userRepository, times(0)).delete(any());
    }

    @Test
    void verifyAtLeastAtMost() {
        verify(userRepository, atLeast(1)).save(any());
        verify(userRepository, atLeastOnce()).save(any());
        verify(userRepository, atMost(5)).findById(anyLong());
    }

    @Test
    void verifyNoMoreInteractions() {
        userService.createUser(new User("John"));

        verify(userRepository).save(any());
        verify(emailService).sendWelcomeEmail(any());
        
        // Verify no other methods were called on these mocks
        verifyNoMoreInteractions(userRepository, emailService);
    }

    @Test
    void verifyZeroInteractions() {
        // Verify mock was never used
        verifyNoInteractions(emailService);
    }

    @Test
    void verifyOrder() {
        userService.createAndNotify(new User("John"));

        // Verify order of calls
        InOrder inOrder = inOrder(userRepository, emailService);
        inOrder.verify(userRepository).save(any());
        inOrder.verify(emailService).sendWelcomeEmail(any());
    }

    @Test
    void verifyWithArgumentCaptor() {
        userService.createUser(new User("John", "john@email.com"));

        // Capture the argument passed to save()
        ArgumentCaptor<User> userCaptor = ArgumentCaptor.forClass(User.class);
        verify(userRepository).save(userCaptor.capture());

        // Assert on captured value
        User captured = userCaptor.getValue();
        assertEquals("John", captured.getName());
        assertEquals("john@email.com", captured.getEmail());
    }

    @Test
    void captureMultipleInvocations() {
        userService.createUsers(List.of(
            new User("John"),
            new User("Jane")
        ));

        ArgumentCaptor<User> captor = ArgumentCaptor.forClass(User.class);
        verify(userRepository, times(2)).save(captor.capture());

        List<User> capturedUsers = captor.getAllValues();
        assertEquals(2, capturedUsers.size());
        assertEquals("John", capturedUsers.get(0).getName());
        assertEquals("Jane", capturedUsers.get(1).getName());
    }
}
```

### 2.7 Spying (Partial Mocks)

```java
class SpyDemo {

    @Spy
    private ArrayList<String> spyList = new ArrayList<>();

    @Spy
    private UserService spyService = new UserService();

    @Test
    void spyBasics() {
        // Spy wraps a real object
        // Real methods are called unless stubbed
        
        spyList.add("one");
        spyList.add("two");

        // Real method was called
        assertEquals(2, spyList.size());
        assertTrue(spyList.contains("one"));

        // Verify interactions
        verify(spyList).add("one");
        verify(spyList).add("two");
    }

    @Test
    void partialStubbing() {
        // Stub specific method, others use real implementation
        doReturn(100).when(spyList).size();
        
        spyList.add("one");
        assertEquals(100, spyList.size()); // Stubbed
        assertTrue(spyList.contains("one")); // Real method
    }

    @Test
    void spyVsMock() {
        // MOCK: All methods return default values unless stubbed
        List<String> mockList = mock(ArrayList.class);
        mockList.add("one");
        assertEquals(0, mockList.size()); // Returns 0 (default for int)

        // SPY: Real methods called unless stubbed
        List<String> spyList = spy(new ArrayList<>());
        spyList.add("one");
        assertEquals(1, spyList.size()); // Real size() called
    }

    @Test
    void whenToUseSpy() {
        // Use spy when you want to test real behavior
        // but need to stub specific methods (e.g., external calls)
        
        UserService realService = new UserService();
        UserService spyService = spy(realService);

        // Stub only the external call
        doReturn(true).when(spyService).sendNotification(any());

        // Rest of the logic uses real implementation
        spyService.processUser(new User());
    }
}
```

### 2.8 BDD Style with Mockito

```java
import static org.mockito.BDDMockito.*;

class BDDStyleDemo {

    @Mock
    private UserRepository userRepository;

    @Test
    void bddStyleTest() {
        // Given
        User user = new User("John", "john@email.com");
        given(userRepository.findByEmail("john@email.com"))
            .willReturn(Optional.of(user));

        // When
        Optional<User> result = userRepository.findByEmail("john@email.com");

        // Then
        then(userRepository).should().findByEmail("john@email.com");
        assertThat(result).isPresent();
        assertThat(result.get().getName()).isEqualTo("John");
    }

    @Test
    void bddVoidMethods() {
        // Given
        willDoNothing().given(userRepository).delete(any());
        
        // Given - throw exception
        willThrow(new RuntimeException("Error"))
            .given(userRepository).delete(999L);
    }

    @Test
    void bddVerification() {
        userService.processUser(1L);

        // Then - verify
        then(userRepository).should(times(1)).findById(1L);
        then(emailService).should(never()).sendEmail(any());
        then(userRepository).shouldHaveNoMoreInteractions();
    }
}
```

---

## 3. Spring Boot Test Integration

### 3.1 Testing Annotations Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    Spring Boot Test Annotations                                  │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐   │
│  │                     @SpringBootTest                                      │   │
│  │  • Loads full application context                                       │   │
│  │  • Slowest but most comprehensive                                       │   │
│  │  • Use for integration tests                                            │   │
│  └─────────────────────────────────────────────────────────────────────────┘   │
│                                    │                                             │
│  ┌──────────────────┬──────────────┴──────────────┬──────────────────┐          │
│  │                  │                             │                  │          │
│  ▼                  ▼                             ▼                  ▼          │
│  ┌────────────┐  ┌────────────┐  ┌────────────────┐  ┌────────────────┐        │
│  │@WebMvcTest │  │@DataJpaTest│  │@WebFluxTest    │  │@JsonTest       │        │
│  │            │  │            │  │                │  │                │        │
│  │ Controllers│  │ Repository │  │ Reactive       │  │ JSON           │        │
│  │ only       │  │ only       │  │ Controllers    │  │ serialization  │        │
│  │ MockMvc    │  │ Embedded DB│  │ WebTestClient  │  │                │        │
│  └────────────┘  └────────────┘  └────────────────┘  └────────────────┘        │
│                                                                                  │
│  OTHER SLICE TESTS:                                                             │
│  @RestClientTest    - REST client testing                                       │
│  @JdbcTest          - JDBC without full JPA                                     │
│  @DataMongoTest     - MongoDB repository testing                                │
│  @DataRedisTest     - Redis testing                                             │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 @SpringBootTest Integration Tests

```java
@SpringBootTest
@AutoConfigureMockMvc
class OrderServiceIntegrationTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private OrderRepository orderRepository;

    @MockBean  // Replaces bean in context with mock
    private PaymentGateway paymentGateway;

    @SpyBean   // Wraps existing bean with spy
    private EmailService emailService;

    @BeforeEach
    void setup() {
        orderRepository.deleteAll();
    }

    @Test
    void createOrder_Success() throws Exception {
        // Given
        when(paymentGateway.charge(any())).thenReturn(PaymentResult.success());

        // When
        mockMvc.perform(post("/api/orders")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""
                    {
                        "userId": 1,
                        "items": [{"productId": 1, "quantity": 2}]
                    }
                    """))
            .andExpect(status().isCreated())
            .andExpect(jsonPath("$.status").value("CREATED"));

        // Then
        assertEquals(1, orderRepository.count());
        verify(emailService).sendOrderConfirmation(any());
    }
}
```

### 3.3 @WebMvcTest - Controller Testing

```java
@WebMvcTest(UserController.class)
class UserControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private UserService userService;

    @Test
    void getUser_Success() throws Exception {
        // Given
        User user = new User(1L, "John", "john@email.com");
        when(userService.findById(1L)).thenReturn(Optional.of(user));

        // When & Then
        mockMvc.perform(get("/api/users/1"))
            .andExpect(status().isOk())
            .andExpect(content().contentType(MediaType.APPLICATION_JSON))
            .andExpect(jsonPath("$.id").value(1))
            .andExpect(jsonPath("$.name").value("John"))
            .andExpect(jsonPath("$.email").value("john@email.com"));
    }

    @Test
    void getUser_NotFound() throws Exception {
        when(userService.findById(999L)).thenReturn(Optional.empty());

        mockMvc.perform(get("/api/users/999"))
            .andExpect(status().isNotFound());
    }

    @Test
    void createUser_ValidationError() throws Exception {
        mockMvc.perform(post("/api/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""
                    {"name": "", "email": "invalid"}
                    """))
            .andExpect(status().isBadRequest())
            .andExpect(jsonPath("$.errors").isArray());
    }
}
```

### 3.4 @DataJpaTest - Repository Testing

```java
@DataJpaTest
@AutoConfigureTestDatabase(replace = Replace.NONE)  // Use real DB
class UserRepositoryTest {

    @Autowired
    private TestEntityManager entityManager;

    @Autowired
    private UserRepository userRepository;

    @BeforeEach
    void setup() {
        User user1 = new User("John", "john@email.com", 25);
        User user2 = new User("Jane", "jane@email.com", 30);
        entityManager.persist(user1);
        entityManager.persist(user2);
        entityManager.flush();
    }

    @Test
    void findByEmail_Success() {
        Optional<User> found = userRepository.findByEmail("john@email.com");
        
        assertThat(found).isPresent();
        assertThat(found.get().getName()).isEqualTo("John");
    }

    @Test
    void findByAgeGreaterThan() {
        List<User> users = userRepository.findByAgeGreaterThan(27);
        
        assertThat(users).hasSize(1);
        assertThat(users.get(0).getName()).isEqualTo("Jane");
    }

    @Test
    void customQuery() {
        List<User> adults = userRepository.findAdults();  // Custom @Query
        
        assertThat(adults).hasSize(2);
    }
}
```

### 3.5 @MockBean vs @Mock

```java
// @MockBean: Replaces Spring bean in application context
// Use when: Testing Spring-managed components

@SpringBootTest
class WithMockBean {
    @MockBean  // Replaces UserService bean in Spring context
    private UserService userService;
    
    @Autowired
    private UserController controller;  // Gets injected with mock
}

// @Mock: Creates mock without Spring context
// Use when: Unit testing without Spring

@ExtendWith(MockitoExtension.class)
class WithMock {
    @Mock  // Just a mock, no Spring involved
    private UserService userService;
    
    @InjectMocks
    private UserController controller;  // Manual injection
}

// Performance comparison:
// @Mock + @InjectMocks: ~10ms per test
// @MockBean + @SpringBootTest: ~2-5 seconds per test
```

### 3.6 Test Configuration

```java
@SpringBootTest
@TestPropertySource(properties = {
    "spring.datasource.url=jdbc:h2:mem:testdb",
    "payment.gateway.url=http://mock-payment"
})
@ActiveProfiles("test")
class ConfiguredTest {
    // Uses test properties and profile
}

// Or use a test configuration class
@TestConfiguration
class TestConfig {

    @Bean
    @Primary  // Override the real bean
    public PaymentGateway mockPaymentGateway() {
        PaymentGateway mock = mock(PaymentGateway.class);
        when(mock.charge(any())).thenReturn(PaymentResult.success());
        return mock;
    }
}

@SpringBootTest
@Import(TestConfig.class)
class WithTestConfig {
    @Autowired
    private PaymentGateway paymentGateway;  // Gets test mock
}
```

---

## 4. MockMvc & REST API Testing

### 4.1 MockMvc Setup and Basics

```java
@WebMvcTest(OrderController.class)
class OrderControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private OrderService orderService;

    // GET Request
    @Test
    void getOrders() throws Exception {
        when(orderService.findAll())
            .thenReturn(List.of(
                new Order(1L, "PENDING"),
                new Order(2L, "COMPLETED")
            ));

        mockMvc.perform(get("/api/orders"))
            .andExpect(status().isOk())
            .andExpect(content().contentType(MediaType.APPLICATION_JSON))
            .andExpect(jsonPath("$", hasSize(2)))
            .andExpect(jsonPath("$[0].id").value(1))
            .andExpect(jsonPath("$[0].status").value("PENDING"))
            .andExpect(jsonPath("$[1].status").value("COMPLETED"));
    }

    // GET with path variable
    @Test
    void getOrderById() throws Exception {
        when(orderService.findById(1L))
            .thenReturn(Optional.of(new Order(1L, "PENDING")));

        mockMvc.perform(get("/api/orders/{id}", 1))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.id").value(1));
    }

    // GET with query parameters
    @Test
    void searchOrders() throws Exception {
        mockMvc.perform(get("/api/orders/search")
                .param("status", "PENDING")
                .param("page", "0")
                .param("size", "10"))
            .andExpect(status().isOk());
    }

    // POST Request
    @Test
    void createOrder() throws Exception {
        Order created = new Order(1L, "CREATED");
        when(orderService.create(any())).thenReturn(created);

        mockMvc.perform(post("/api/orders")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""
                    {
                        "userId": 1,
                        "items": [
                            {"productId": 100, "quantity": 2},
                            {"productId": 101, "quantity": 1}
                        ]
                    }
                    """))
            .andExpect(status().isCreated())
            .andExpect(header().exists("Location"))
            .andExpect(jsonPath("$.id").value(1))
            .andExpect(jsonPath("$.status").value("CREATED"));
    }

    // PUT Request
    @Test
    void updateOrder() throws Exception {
        Order updated = new Order(1L, "PROCESSING");
        when(orderService.update(eq(1L), any())).thenReturn(updated);

        mockMvc.perform(put("/api/orders/{id}", 1)
                .contentType(MediaType.APPLICATION_JSON)
                .content("""
                    {"status": "PROCESSING"}
                    """))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.status").value("PROCESSING"));
    }

    // DELETE Request
    @Test
    void deleteOrder() throws Exception {
        doNothing().when(orderService).delete(1L);

        mockMvc.perform(delete("/api/orders/{id}", 1))
            .andExpect(status().isNoContent());

        verify(orderService).delete(1L);
    }

    // PATCH Request
    @Test
    void patchOrder() throws Exception {
        mockMvc.perform(patch("/api/orders/{id}", 1)
                .contentType(MediaType.APPLICATION_JSON)
                .content("""
                    {"status": "SHIPPED"}
                    """))
            .andExpect(status().isOk());
    }
}
```

### 4.2 Request Customization

```java
class RequestCustomizationTest {

    @Test
    void requestWithHeaders() throws Exception {
        mockMvc.perform(get("/api/orders")
                .header("Authorization", "Bearer token123")
                .header("X-Request-Id", "req-001")
                .accept(MediaType.APPLICATION_JSON))
            .andExpect(status().isOk());
    }

    @Test
    void requestWithCookies() throws Exception {
        mockMvc.perform(get("/api/orders")
                .cookie(new Cookie("session", "abc123")))
            .andExpect(status().isOk());
    }

    @Test
    void multipartUpload() throws Exception {
        MockMultipartFile file = new MockMultipartFile(
            "file",                    // Parameter name
            "test.csv",                // Original filename
            MediaType.TEXT_PLAIN_VALUE,
            "id,name\n1,John".getBytes()
        );

        mockMvc.perform(multipart("/api/import")
                .file(file)
                .param("type", "users"))
            .andExpect(status().isOk());
    }

    @Test
    void formSubmission() throws Exception {
        mockMvc.perform(post("/api/login")
                .contentType(MediaType.APPLICATION_FORM_URLENCODED)
                .param("username", "john")
                .param("password", "secret"))
            .andExpect(status().is3xxRedirection());
    }
}
```

### 4.3 Response Assertions

```java
class ResponseAssertionsTest {

    @Test
    void assertStatus() throws Exception {
        mockMvc.perform(get("/api/orders"))
            .andExpect(status().isOk())              // 200
            .andExpect(status().is2xxSuccessful());  // 2xx range

        mockMvc.perform(post("/api/orders"))
            .andExpect(status().isCreated());        // 201

        mockMvc.perform(get("/api/notfound"))
            .andExpect(status().isNotFound());       // 404

        mockMvc.perform(get("/api/error"))
            .andExpect(status().is5xxServerError()); // 5xx range
    }

    @Test
    void assertHeaders() throws Exception {
        mockMvc.perform(post("/api/orders"))
            .andExpect(header().exists("Location"))
            .andExpect(header().string("Location", containsString("/api/orders/")))
            .andExpect(header().string("Content-Type", MediaType.APPLICATION_JSON_VALUE));
    }

    @Test
    void assertJsonPath() throws Exception {
        mockMvc.perform(get("/api/orders/1"))
            // Value assertions
            .andExpect(jsonPath("$.id").value(1))
            .andExpect(jsonPath("$.status").value("PENDING"))
            
            // Type assertions
            .andExpect(jsonPath("$.id").isNumber())
            .andExpect(jsonPath("$.items").isArray())
            .andExpect(jsonPath("$.user").isMap())
            
            // Collection assertions
            .andExpect(jsonPath("$.items", hasSize(3)))
            .andExpect(jsonPath("$.items[0].name").value("Product A"))
            
            // Existence assertions
            .andExpect(jsonPath("$.createdAt").exists())
            .andExpect(jsonPath("$.deletedAt").doesNotExist())
            
            // Matchers
            .andExpect(jsonPath("$.total").value(greaterThan(0.0)))
            .andExpect(jsonPath("$.status").value(oneOf("PENDING", "PROCESSING")));
    }

    @Test
    void assertContent() throws Exception {
        mockMvc.perform(get("/api/orders"))
            .andExpect(content().contentType(MediaType.APPLICATION_JSON))
            .andExpect(content().json("""
                [
                    {"id": 1, "status": "PENDING"},
                    {"id": 2, "status": "COMPLETED"}
                ]
                """))
            .andExpect(content().string(containsString("PENDING")));
    }

    @Test
    void extractResponse() throws Exception {
        MvcResult result = mockMvc.perform(get("/api/orders/1"))
            .andExpect(status().isOk())
            .andReturn();

        String content = result.getResponse().getContentAsString();
        Order order = objectMapper.readValue(content, Order.class);
        
        assertEquals(1L, order.getId());
    }
}
```

### 4.4 Validation Testing

```java
@WebMvcTest(UserController.class)
class ValidationTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private UserService userService;

    @Test
    void createUser_ValidationErrors() throws Exception {
        mockMvc.perform(post("/api/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content("""
                    {
                        "name": "",
                        "email": "invalid-email",
                        "age": -5
                    }
                    """))
            .andExpect(status().isBadRequest())
            .andExpect(jsonPath("$.errors", hasSize(3)))
            .andExpect(jsonPath("$.errors[*].field", 
                containsInAnyOrder("name", "email", "age")));
    }

    @Test
    void createUser_MissingRequiredFields() throws Exception {
        mockMvc.perform(post("/api/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content("{}"))
            .andExpect(status().isBadRequest())
            .andExpect(jsonPath("$.errors[*].message", 
                hasItem(containsString("must not be blank"))));
    }
}
```

### 4.5 Security Testing

```java
@WebMvcTest(OrderController.class)
@Import(SecurityConfig.class)
class SecurityTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private OrderService orderService;

    @Test
    void unauthenticated_Unauthorized() throws Exception {
        mockMvc.perform(get("/api/orders"))
            .andExpect(status().isUnauthorized());
    }

    @Test
    @WithMockUser(username = "john", roles = {"USER"})
    void authenticatedUser_Success() throws Exception {
        when(orderService.findAll()).thenReturn(List.of());

        mockMvc.perform(get("/api/orders"))
            .andExpect(status().isOk());
    }

    @Test
    @WithMockUser(username = "john", roles = {"USER"})
    void userAccessingAdminEndpoint_Forbidden() throws Exception {
        mockMvc.perform(delete("/api/admin/orders/1"))
            .andExpect(status().isForbidden());
    }

    @Test
    @WithMockUser(username = "admin", roles = {"ADMIN"})
    void adminAccess_Success() throws Exception {
        doNothing().when(orderService).delete(1L);

        mockMvc.perform(delete("/api/admin/orders/1"))
            .andExpect(status().isNoContent());
    }

    @Test
    void withJwtToken() throws Exception {
        String token = generateTestJwt("john", "USER");

        mockMvc.perform(get("/api/orders")
                .header("Authorization", "Bearer " + token))
            .andExpect(status().isOk());
    }
}
```

---

## 5. Database Testing with Testcontainers

### 5.1 Testcontainers Setup

```xml
<!-- pom.xml -->
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>testcontainers</artifactId>
    <scope>test</scope>
</dependency>
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>junit-jupiter</artifactId>
    <scope>test</scope>
</dependency>
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>postgresql</artifactId>
    <scope>test</scope>
</dependency>
```

### 5.2 Basic Testcontainers Usage

```java
@SpringBootTest
@Testcontainers
class OrderRepositoryIntegrationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15-alpine")
        .withDatabaseName("testdb")
        .withUsername("test")
        .withPassword("test");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }

    @Autowired
    private OrderRepository orderRepository;

    @BeforeEach
    void setup() {
        orderRepository.deleteAll();
    }

    @Test
    void saveAndFindOrder() {
        Order order = new Order("user-1", OrderStatus.PENDING);
        Order saved = orderRepository.save(order);

        Optional<Order> found = orderRepository.findById(saved.getId());

        assertThat(found).isPresent();
        assertThat(found.get().getUserId()).isEqualTo("user-1");
    }

    @Test
    void findByStatus() {
        orderRepository.save(new Order("user-1", OrderStatus.PENDING));
        orderRepository.save(new Order("user-2", OrderStatus.COMPLETED));
        orderRepository.save(new Order("user-3", OrderStatus.PENDING));

        List<Order> pendingOrders = orderRepository.findByStatus(OrderStatus.PENDING);

        assertThat(pendingOrders).hasSize(2);
    }
}
```

### 5.3 Shared Container (Performance Optimization)

```java
// Base class with shared container
@Testcontainers
public abstract class AbstractIntegrationTest {

    @Container
    protected static final PostgreSQLContainer<?> postgres;

    static {
        postgres = new PostgreSQLContainer<>("postgres:15-alpine")
            .withDatabaseName("testdb")
            .withUsername("test")
            .withPassword("test")
            .withReuse(true);  // Reuse container across test runs
        postgres.start();
    }

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }
}

// Test classes extend base
@SpringBootTest
class OrderServiceTest extends AbstractIntegrationTest {

    @Autowired
    private OrderService orderService;

    @Test
    void createOrder() {
        // Uses shared PostgreSQL container
    }
}

@SpringBootTest
class UserServiceTest extends AbstractIntegrationTest {

    @Autowired
    private UserService userService;

    @Test
    void createUser() {
        // Uses same PostgreSQL container
    }
}
```

### 5.4 Multiple Containers

```java
@SpringBootTest
@Testcontainers
class FullIntegrationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15");

    @Container
    static GenericContainer<?> redis = new GenericContainer<>("redis:7-alpine")
        .withExposedPorts(6379);

    @Container
    static KafkaContainer kafka = new KafkaContainer(
        DockerImageName.parse("confluentinc/cp-kafka:7.5.0"));

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        // PostgreSQL
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);

        // Redis
        registry.add("spring.data.redis.host", redis::getHost);
        registry.add("spring.data.redis.port", () -> redis.getMappedPort(6379));

        // Kafka
        registry.add("spring.kafka.bootstrap-servers", kafka::getBootstrapServers);
    }

    @Test
    void fullIntegrationTest() {
        // Test with all infrastructure
    }
}
```

### 5.5 Database Migration Testing

```java
@SpringBootTest
@Testcontainers
class FlywayMigrationTest {

    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15");

    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
        registry.add("spring.flyway.enabled", () -> true);
    }

    @Autowired
    private Flyway flyway;

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @Test
    void migrationsAreApplied() {
        // Verify migrations ran
        MigrationInfo[] appliedMigrations = flyway.info().applied();
        assertThat(appliedMigrations).isNotEmpty();
    }

    @Test
    void tablesExist() {
        // Verify expected tables
        String sql = """
            SELECT table_name FROM information_schema.tables 
            WHERE table_schema = 'public'
            """;
        
        List<String> tables = jdbcTemplate.queryForList(sql, String.class);
        
        assertThat(tables).contains("users", "orders", "order_items");
    }
}
```

---

## 6. Advanced Mocking Patterns

### 6.1 Static Method Mocking

```java
import org.mockito.MockedStatic;

class StaticMockingDemo {

    @Test
    void mockStaticMethod() {
        try (MockedStatic<UUID> mockedUUID = mockStatic(UUID.class)) {
            UUID fixedUUID = UUID.fromString("123e4567-e89b-12d3-a456-426614174000");
            mockedUUID.when(UUID::randomUUID).thenReturn(fixedUUID);

            // Now UUID.randomUUID() returns our fixed value
            Order order = orderService.createOrder(request);
            
            assertEquals("123e4567-e89b-12d3-a456-426614174000", 
                        order.getOrderId());
        }
        // Static mock is automatically reset after try block
    }

    @Test
    void mockStaticWithArguments() {
        try (MockedStatic<LocalDateTime> mockedTime = mockStatic(LocalDateTime.class)) {
            LocalDateTime fixedTime = LocalDateTime.of(2026, 2, 28, 10, 0);
            
            mockedTime.when(LocalDateTime::now).thenReturn(fixedTime);
            mockedTime.when(() -> LocalDateTime.parse(anyString()))
                .thenReturn(fixedTime);

            // Test time-dependent logic
            Order order = orderService.createOrder(request);
            
            assertEquals(fixedTime, order.getCreatedAt());
        }
    }
}
```

### 6.2 Constructor Mocking

```java
import org.mockito.MockedConstruction;

class ConstructorMockingDemo {

    @Test
    void mockConstructor() {
        try (MockedConstruction<PaymentProcessor> mocked = 
                mockConstruction(PaymentProcessor.class, (mock, context) -> {
                    when(mock.process(any())).thenReturn(PaymentResult.success());
                })) {
            
            // Any new PaymentProcessor() inside orderService will be mocked
            Order order = orderService.createOrder(request);
            
            assertEquals(OrderStatus.PAID, order.getStatus());
            
            // Verify mock was created
            assertEquals(1, mocked.constructed().size());
        }
    }

    @Test
    void mockConstructorWithArguments() {
        try (MockedConstruction<DatabaseConnection> mocked = 
                mockConstruction(DatabaseConnection.class, (mock, context) -> {
                    // Access constructor arguments
                    List<?> args = context.arguments();
                    String connectionString = (String) args.get(0);
                    
                    when(mock.isConnected()).thenReturn(true);
                })) {
            
            dataService.connect("jdbc:postgresql://localhost/test");
        }
    }
}
```

### 6.3 Mocking Final Classes and Methods

```java
// Enable mock-maker-inline in src/test/resources/mockito-extensions
// org.mockito.plugins.MockMaker file containing: mock-maker-inline

class FinalClassMockingDemo {

    @Mock
    private FinalService finalService;  // Final class can be mocked

    @Test
    void mockFinalMethod() {
        when(finalService.finalMethod()).thenReturn("mocked");
        
        assertEquals("mocked", finalService.finalMethod());
    }
}

// Alternative: Using mockito-inline dependency
// <artifactId>mockito-inline</artifactId>
```

### 6.4 Deep Stubs

```java
class DeepStubsDemo {

    @Mock(answer = Answers.RETURNS_DEEP_STUBS)
    private OrderService orderService;

    @Test
    void withDeepStubs() {
        // Without deep stubs, you'd need:
        // Order mockOrder = mock(Order.class);
        // User mockUser = mock(User.class);
        // when(orderService.findById(1L)).thenReturn(mockOrder);
        // when(mockOrder.getUser()).thenReturn(mockUser);
        // when(mockUser.getName()).thenReturn("John");

        // With deep stubs:
        when(orderService.findById(1L).getUser().getName())
            .thenReturn("John");

        assertEquals("John", orderService.findById(1L).getUser().getName());
    }
}
```

### 6.5 Lenient Stubbing

```java
class LenientStubbingDemo {

    @Test
    void strictStubbing() {
        // By default, Mockito is strict - unused stubs cause failure
        when(userService.findById(1L)).thenReturn(user1);
        when(userService.findById(2L)).thenReturn(user2);  // Unused - fails!
        
        userService.findById(1L);  // Only this is used
    }

    @Test
    void lenientStubbing() {
        // Mark specific stub as lenient
        lenient().when(userService.findById(2L)).thenReturn(user2);
        
        // Or use in setup for common stubs
    }
}

@ExtendWith(MockitoExtension.class)
@MockitoSettings(strictness = Strictness.LENIENT)  // All stubs lenient
class LenientTestClass {
    // All stubbing in this class is lenient
}
```

### 6.6 Custom Answer Implementations

```java
class CustomAnswerDemo {

    @Test
    void customAnswer() {
        when(userRepository.findById(anyLong())).thenAnswer(new Answer<Optional<User>>() {
            private int callCount = 0;

            @Override
            public Optional<User> answer(InvocationOnMock invocation) {
                callCount++;
                Long id = invocation.getArgument(0);
                
                if (callCount == 1) {
                    return Optional.empty();  // First call returns empty
                }
                return Optional.of(new User(id, "User " + id));
            }
        });
    }

    @Test
    void delayedAnswer() {
        when(externalService.call(any())).thenAnswer(invocation -> {
            Thread.sleep(100);  // Simulate network delay
            return new Response("OK");
        });
    }

    @Test
    void callRealMethodAnswer() {
        UserService spy = spy(new UserService());
        
        when(spy.processUser(any()))
            .thenAnswer(AdditionalAnswers.answersWithDelay(100, 
                invocation -> invocation.callRealMethod()));
    }
}
```

---

## 7. Testing Best Practices & Patterns

### 7.1 Test Naming Conventions

```java
class NamingConventionsDemo {

    // Pattern 1: methodName_scenario_expectedResult
    @Test
    void createUser_validInput_returnsCreatedUser() { }
    
    @Test
    void createUser_nullEmail_throwsValidationException() { }

    // Pattern 2: should_expectedBehavior_when_scenario
    @Test
    void should_returnUser_when_validIdProvided() { }
    
    @Test
    void should_throwException_when_userNotFound() { }

    // Pattern 3: given_when_then (BDD style)
    @Test
    void givenValidUser_whenCreate_thenReturnSavedUser() { }

    // Pattern 4: Using @DisplayName for readability
    @Test
    @DisplayName("Creating user with valid data should return saved user")
    void createUserValid() { }
}
```

### 7.2 Test Structure - AAA Pattern

```java
class AAAPatternDemo {

    @Test
    void createOrder_Success() {
        // Arrange - Setup test data and mocks
        User user = new User(1L, "John", "john@email.com");
        List<OrderItem> items = List.of(
            new OrderItem("SKU001", 2),
            new OrderItem("SKU002", 1)
        );
        CreateOrderRequest request = new CreateOrderRequest(user.getId(), items);
        
        when(userRepository.findById(1L)).thenReturn(Optional.of(user));
        when(inventoryService.checkAvailability(anyList())).thenReturn(true);
        when(orderRepository.save(any())).thenAnswer(inv -> {
            Order order = inv.getArgument(0);
            order.setId(100L);
            return order;
        });

        // Act - Execute the method under test
        Order result = orderService.createOrder(request);

        // Assert - Verify the results
        assertThat(result.getId()).isEqualTo(100L);
        assertThat(result.getStatus()).isEqualTo(OrderStatus.CREATED);
        assertThat(result.getItems()).hasSize(2);
        
        verify(orderRepository).save(any(Order.class));
        verify(inventoryService).reserveItems(anyList());
    }
}
```

### 7.3 Test Data Builders

```java
// Builder pattern for test data
public class UserBuilder {
    private Long id = 1L;
    private String name = "John Doe";
    private String email = "john@example.com";
    private int age = 30;
    private UserStatus status = UserStatus.ACTIVE;

    public static UserBuilder aUser() {
        return new UserBuilder();
    }

    public UserBuilder withId(Long id) {
        this.id = id;
        return this;
    }

    public UserBuilder withName(String name) {
        this.name = name;
        return this;
    }

    public UserBuilder withEmail(String email) {
        this.email = email;
        return this;
    }

    public UserBuilder inactive() {
        this.status = UserStatus.INACTIVE;
        return this;
    }

    public User build() {
        return new User(id, name, email, age, status);
    }
}

// Usage in tests
class UserServiceTest {

    @Test
    void processActiveUser() {
        User user = aUser()
            .withName("Jane")
            .withEmail("jane@example.com")
            .build();

        // Test with active user (default)
    }

    @Test
    void skipInactiveUser() {
        User user = aUser()
            .inactive()
            .build();

        // Test with inactive user
    }
}
```

### 7.4 Test Fixtures

```java
class OrderServiceTest {

    // Shared test fixtures
    private static User testUser;
    private static List<Product> testProducts;

    @BeforeAll
    static void setupFixtures() {
        testUser = aUser().withId(1L).build();
        testProducts = List.of(
            aProduct().withSku("SKU001").withPrice(10.0).build(),
            aProduct().withSku("SKU002").withPrice(20.0).build()
        );
    }

    @BeforeEach
    void setupMocks() {
        when(userRepository.findById(1L)).thenReturn(Optional.of(testUser));
        testProducts.forEach(p -> 
            when(productRepository.findBySku(p.getSku()))
                .thenReturn(Optional.of(p)));
    }

    @Test
    void createOrderWithFixtures() {
        // Tests use pre-configured fixtures
    }
}
```

### 7.5 Testing Exceptions

```java
class ExceptionTestingDemo {

    @Test
    void assertThrowsWithType() {
        assertThrows(UserNotFoundException.class, 
            () -> userService.findById(999L));
    }

    @Test
    void assertThrowsWithMessage() {
        UserNotFoundException exception = assertThrows(
            UserNotFoundException.class,
            () -> userService.findById(999L)
        );
        
        assertThat(exception.getMessage()).contains("999");
        assertThat(exception.getErrorCode()).isEqualTo("USER_NOT_FOUND");
    }

    @Test
    void assertThrowsWithCause() {
        DataAccessException exception = assertThrows(
            DataAccessException.class,
            () -> orderService.save(order)
        );
        
        assertThat(exception.getCause())
            .isInstanceOf(SQLException.class);
    }

    @Test
    void assertDoesNotThrow() {
        assertDoesNotThrow(() -> userService.validateEmail("valid@email.com"));
    }
}
```

### 7.6 Async Testing

```java
class AsyncTestingDemo {

    @Test
    void testAsyncMethod() throws Exception {
        CompletableFuture<Order> future = orderService.createOrderAsync(request);
        
        Order result = future.get(5, TimeUnit.SECONDS);
        
        assertThat(result.getStatus()).isEqualTo(OrderStatus.CREATED);
    }

    @Test
    void testWithAwaitility() {
        // Trigger async operation
        orderService.processOrderAsync(orderId);
        
        // Wait for condition
        await()
            .atMost(10, TimeUnit.SECONDS)
            .pollInterval(500, TimeUnit.MILLISECONDS)
            .untilAsserted(() -> {
                Order order = orderRepository.findById(orderId).orElseThrow();
                assertThat(order.getStatus()).isEqualTo(OrderStatus.PROCESSED);
            });
    }

    @Test
    void testEventualConsistency() {
        // Publish event
        eventPublisher.publish(new OrderCreatedEvent(orderId));
        
        // Wait for downstream effect
        await()
            .atMost(Duration.ofSeconds(5))
            .until(() -> inventoryRepository.findReservation(orderId).isPresent());
    }
}
```

### 7.7 Test Coverage Guidelines

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                      Test Coverage Guidelines                                    │
│                                                                                  │
│  WHAT TO TEST:                                                                   │
│  ✓ Business logic and calculations                                              │
│  ✓ Edge cases and boundary conditions                                           │
│  ✓ Error handling and exceptions                                                │
│  ✓ Integration points (APIs, databases)                                         │
│  ✓ Security-sensitive code                                                      │
│                                                                                  │
│  WHAT NOT TO TEST:                                                               │
│  ✗ Simple getters/setters                                                       │
│  ✗ Configuration classes                                                        │
│  ✗ Third-party library code                                                     │
│  ✗ Generated code                                                               │
│                                                                                  │
│  COVERAGE TARGETS:                                                               │
│  • Unit Tests: 80%+ line coverage for business logic                            │
│  • Integration Tests: Critical paths and workflows                              │
│  • End-to-End: Happy path + main failure scenarios                              │
│                                                                                  │
│  FOCUS ON:                                                                       │
│  • Branch coverage over line coverage                                           │
│  • Testing behavior, not implementation                                         │
│  • Meaningful assertions over coverage percentage                               │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 8. Interview Questions

### Basic Questions

**Q1: What is the difference between @Mock and @MockBean?**
> - `@Mock` creates a mock without Spring context (Mockito only). Faster, used for unit tests.
> - `@MockBean` replaces a Spring bean in the application context. Used with `@SpringBootTest`.

**Q2: Explain @InjectMocks annotation.**
> `@InjectMocks` creates an instance of the class and injects mocks into it via constructor/setter/field injection. It's useful for unit testing without Spring context.

**Q3: What's the difference between mock() and spy()?**
> - `mock()`: All methods return default values unless stubbed. No real methods called.
> - `spy()`: Wraps a real object. Real methods called unless stubbed. Useful for partial mocking.

**Q4: How do you verify method invocations in Mockito?**
```java
verify(mock).method();                    // Called exactly once
verify(mock, times(2)).method();          // Called exactly twice
verify(mock, never()).method();           // Never called
verify(mock, atLeast(1)).method();        // Called at least once
verifyNoMoreInteractions(mock);           // No other methods called
```

**Q5: What is ArgumentCaptor and when to use it?**
> ArgumentCaptor captures arguments passed to mock methods for assertion:
```java
ArgumentCaptor<User> captor = ArgumentCaptor.forClass(User.class);
verify(userRepo).save(captor.capture());
assertEquals("John", captor.getValue().getName());
```

### Intermediate Questions

**Q6: Explain the difference between @BeforeEach and @BeforeAll.**
> - `@BeforeEach`: Runs before each test method. Good for resetting state.
> - `@BeforeAll`: Runs once before all tests in class (must be static unless `@TestInstance(PER_CLASS)`). Good for expensive setup.

**Q7: How do you test exceptions in JUnit 5?**
```java
// Assert exception is thrown
Exception ex = assertThrows(IllegalArgumentException.class, 
    () -> service.process(null));

// Assert message
assertTrue(ex.getMessage().contains("null"));

// Assert no exception
assertDoesNotThrow(() -> service.process(validInput));
```

**Q8: What is @ParameterizedTest and its sources?**
> Runs same test with different inputs. Sources include:
> - `@ValueSource`: Primitive values
> - `@EnumSource`: Enum values
> - `@CsvSource`: CSV inline data
> - `@CsvFileSource`: External CSV file
> - `@MethodSource`: Custom method providing arguments

**Q9: How do you test private methods?**
> - Don't test private methods directly - test through public methods
> - If needed, use reflection (last resort)
> - Consider if method should be package-private or in separate class
> - High need to test private method = possible design smell

**Q10: Explain Testcontainers and its benefits.**
> Testcontainers provides throwaway instances of databases, message brokers, etc. in Docker containers:
> - Real database instead of H2/mock
> - Consistent with production
> - Isolated per test run
> - Supports PostgreSQL, MySQL, MongoDB, Kafka, Redis, etc.

### Advanced Questions

**Q11: How do you test async methods?**
```java
// Using CompletableFuture
CompletableFuture<Result> future = service.asyncMethod();
Result result = future.get(5, TimeUnit.SECONDS);

// Using Awaitility
await()
    .atMost(10, SECONDS)
    .until(() -> repository.findById(id).isPresent());
```

**Q12: How do you mock static methods?**
```java
try (MockedStatic<UUID> mocked = mockStatic(UUID.class)) {
    mocked.when(UUID::randomUUID).thenReturn(fixedUUID);
    // Test code
}  // Auto-reset after try block
```

**Q13: What is @WebMvcTest and what does it load?**
> `@WebMvcTest` loads only web layer components:
> - Controllers specified in annotation
> - `@ControllerAdvice`, `@JsonComponent`
> - Security filters, WebMvc configurer
> - Does NOT load `@Service`, `@Repository`, `@Component` (need `@MockBean`)

**Q14: How do you test Spring Security with MockMvc?**
```java
// Using @WithMockUser
@Test
@WithMockUser(username = "admin", roles = {"ADMIN"})
void adminAccess() throws Exception {
    mockMvc.perform(get("/admin"))
        .andExpect(status().isOk());
}

// Using SecurityMockMvcRequestPostProcessors
mockMvc.perform(get("/api")
    .with(jwt().authorities(new SimpleGrantedAuthority("ROLE_USER"))))
    .andExpect(status().isOk());
```

**Q15: Design a testing strategy for a microservices application.**
```
1. Unit Tests (70%)
   - Service layer with mocked dependencies
   - Use @Mock, @InjectMocks
   
2. Integration Tests (20%)
   - Repository with Testcontainers
   - Controller with @WebMvcTest
   
3. Contract Tests (5%)
   - Spring Cloud Contract / Pact
   - Producer verifies contracts
   - Consumer tests against stubs
   
4. End-to-End Tests (5%)
   - Full service deployment
   - Docker Compose setup
   - Critical happy paths only
```

---

## 9. Complete CRUD Testing Example

> A comprehensive, production-ready example demonstrating how to test all CRUD operations across all layers (Entity, Repository, Service, Controller).

### 9.1 Project Structure

```
src/
├── main/java/com/example/product/
│   ├── controller/
│   │   └── ProductController.java
│   ├── dto/
│   │   ├── ProductRequest.java
│   │   └── ProductResponse.java
│   ├── entity/
│   │   └── Product.java
│   ├── exception/
│   │   ├── GlobalExceptionHandler.java
│   │   └── ProductNotFoundException.java
│   ├── repository/
│   │   └── ProductRepository.java
│   └── service/
│       ├── ProductService.java
│       └── ProductServiceImpl.java
└── test/java/com/example/product/
    ├── controller/
    │   └── ProductControllerTest.java
    ├── repository/
    │   └── ProductRepositoryTest.java
    └── service/
        └── ProductServiceTest.java
```

### 9.2 Entity Class

```java
package com.example.product.entity;

import jakarta.persistence.*;
import lombok.*;
import java.math.BigDecimal;
import java.time.LocalDateTime;

@Entity
@Table(name = "products")
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Product {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, length = 100)
    private String name;

    @Column(length = 500)
    private String description;

    @Column(nullable = false, precision = 10, scale = 2)
    private BigDecimal price;

    @Column(nullable = false)
    private Integer quantity;

    @Column(name = "category")
    private String category;

    @Column(name = "sku", unique = true, nullable = false)
    private String sku;

    @Column(name = "active")
    private Boolean active = true;

    @Column(name = "created_at", updatable = false)
    private LocalDateTime createdAt;

    @Column(name = "updated_at")
    private LocalDateTime updatedAt;

    @PrePersist
    protected void onCreate() {
        createdAt = LocalDateTime.now();
        updatedAt = LocalDateTime.now();
    }

    @PreUpdate
    protected void onUpdate() {
        updatedAt = LocalDateTime.now();
    }
}
```

### 9.3 DTOs (Request & Response)

```java
package com.example.product.dto;

import jakarta.validation.constraints.*;
import lombok.*;
import java.math.BigDecimal;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class ProductRequest {

    @NotBlank(message = "Product name is required")
    @Size(min = 2, max = 100, message = "Name must be between 2 and 100 characters")
    private String name;

    @Size(max = 500, message = "Description cannot exceed 500 characters")
    private String description;

    @NotNull(message = "Price is required")
    @DecimalMin(value = "0.01", message = "Price must be greater than 0")
    @Digits(integer = 8, fraction = 2, message = "Invalid price format")
    private BigDecimal price;

    @NotNull(message = "Quantity is required")
    @Min(value = 0, message = "Quantity cannot be negative")
    private Integer quantity;

    private String category;

    @NotBlank(message = "SKU is required")
    @Pattern(regexp = "^[A-Z]{3}-[0-9]{6}$", message = "SKU must match pattern XXX-000000")
    private String sku;
}
```

```java
package com.example.product.dto;

import lombok.*;
import java.math.BigDecimal;
import java.time.LocalDateTime;

@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class ProductResponse {
    private Long id;
    private String name;
    private String description;
    private BigDecimal price;
    private Integer quantity;
    private String category;
    private String sku;
    private Boolean active;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}
```

### 9.4 Custom Exception

```java
package com.example.product.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(HttpStatus.NOT_FOUND)
public class ProductNotFoundException extends RuntimeException {
    
    public ProductNotFoundException(Long id) {
        super("Product not found with id: " + id);
    }
    
    public ProductNotFoundException(String sku) {
        super("Product not found with SKU: " + sku);
    }
}
```

```java
package com.example.product.exception;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.util.*;

@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(ProductNotFoundException.class)
    public ResponseEntity<Map<String, Object>> handleProductNotFound(
            ProductNotFoundException ex) {
        Map<String, Object> error = new LinkedHashMap<>();
        error.put("timestamp", LocalDateTime.now());
        error.put("status", HttpStatus.NOT_FOUND.value());
        error.put("error", "Not Found");
        error.put("message", ex.getMessage());
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(error);
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<Map<String, Object>> handleValidationErrors(
            MethodArgumentNotValidException ex) {
        Map<String, Object> error = new LinkedHashMap<>();
        error.put("timestamp", LocalDateTime.now());
        error.put("status", HttpStatus.BAD_REQUEST.value());
        error.put("error", "Validation Failed");
        
        Map<String, String> fieldErrors = new HashMap<>();
        ex.getBindingResult().getAllErrors().forEach(e -> {
            String fieldName = ((FieldError) e).getField();
            String errorMessage = e.getDefaultMessage();
            fieldErrors.put(fieldName, errorMessage);
        });
        error.put("errors", fieldErrors);
        
        return ResponseEntity.badRequest().body(error);
    }
}
```

### 9.5 Repository Interface

```java
package com.example.product.repository;

import com.example.product.entity.Product;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.math.BigDecimal;
import java.util.List;
import java.util.Optional;

@Repository
public interface ProductRepository extends JpaRepository<Product, Long> {

    Optional<Product> findBySku(String sku);

    List<Product> findByCategory(String category);

    List<Product> findByActiveTrue();

    List<Product> findByPriceBetween(BigDecimal minPrice, BigDecimal maxPrice);

    @Query("SELECT p FROM Product p WHERE p.name LIKE %:keyword% OR p.description LIKE %:keyword%")
    List<Product> searchByKeyword(@Param("keyword") String keyword);

    @Query("SELECT p FROM Product p WHERE p.quantity < :threshold AND p.active = true")
    List<Product> findLowStockProducts(@Param("threshold") Integer threshold);

    Page<Product> findByCategoryAndActiveTrue(String category, Pageable pageable);

    @Modifying
    @Query("UPDATE Product p SET p.active = false WHERE p.id = :id")
    int softDeleteById(@Param("id") Long id);

    boolean existsBySku(String sku);
}
```

### 9.6 Service Interface & Implementation

```java
package com.example.product.service;

import com.example.product.dto.ProductRequest;
import com.example.product.dto.ProductResponse;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

import java.util.List;

public interface ProductService {
    
    // CREATE
    ProductResponse createProduct(ProductRequest request);
    
    // READ
    ProductResponse getProductById(Long id);
    ProductResponse getProductBySku(String sku);
    List<ProductResponse> getAllProducts();
    List<ProductResponse> getProductsByCategory(String category);
    Page<ProductResponse> getProducts(Pageable pageable);
    
    // UPDATE
    ProductResponse updateProduct(Long id, ProductRequest request);
    ProductResponse partialUpdateProduct(Long id, ProductRequest request);
    
    // DELETE
    void deleteProduct(Long id);
    void softDeleteProduct(Long id);
}
```

```java
package com.example.product.service;

import com.example.product.dto.ProductRequest;
import com.example.product.dto.ProductResponse;
import com.example.product.entity.Product;
import com.example.product.exception.ProductNotFoundException;
import com.example.product.repository.ProductRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
@Transactional(readOnly = true)
public class ProductServiceImpl implements ProductService {

    private final ProductRepository productRepository;

    // ========== CREATE ==========
    @Override
    @Transactional
    public ProductResponse createProduct(ProductRequest request) {
        log.info("Creating new product with SKU: {}", request.getSku());
        
        // Check if SKU already exists
        if (productRepository.existsBySku(request.getSku())) {
            throw new IllegalArgumentException("Product with SKU " + request.getSku() + " already exists");
        }

        Product product = mapToEntity(request);
        Product savedProduct = productRepository.save(product);
        
        log.info("Product created successfully with ID: {}", savedProduct.getId());
        return mapToResponse(savedProduct);
    }

    // ========== READ ==========
    @Override
    public ProductResponse getProductById(Long id) {
        log.info("Fetching product by ID: {}", id);
        Product product = productRepository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));
        return mapToResponse(product);
    }

    @Override
    public ProductResponse getProductBySku(String sku) {
        log.info("Fetching product by SKU: {}", sku);
        Product product = productRepository.findBySku(sku)
            .orElseThrow(() -> new ProductNotFoundException(sku));
        return mapToResponse(product);
    }

    @Override
    public List<ProductResponse> getAllProducts() {
        log.info("Fetching all products");
        return productRepository.findAll()
            .stream()
            .map(this::mapToResponse)
            .collect(Collectors.toList());
    }

    @Override
    public List<ProductResponse> getProductsByCategory(String category) {
        log.info("Fetching products by category: {}", category);
        return productRepository.findByCategory(category)
            .stream()
            .map(this::mapToResponse)
            .collect(Collectors.toList());
    }

    @Override
    public Page<ProductResponse> getProducts(Pageable pageable) {
        log.info("Fetching products with pagination");
        return productRepository.findAll(pageable)
            .map(this::mapToResponse);
    }

    // ========== UPDATE ==========
    @Override
    @Transactional
    public ProductResponse updateProduct(Long id, ProductRequest request) {
        log.info("Updating product with ID: {}", id);
        
        Product existingProduct = productRepository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));

        // Check SKU uniqueness if changed
        if (!existingProduct.getSku().equals(request.getSku()) 
                && productRepository.existsBySku(request.getSku())) {
            throw new IllegalArgumentException("Product with SKU " + request.getSku() + " already exists");
        }

        existingProduct.setName(request.getName());
        existingProduct.setDescription(request.getDescription());
        existingProduct.setPrice(request.getPrice());
        existingProduct.setQuantity(request.getQuantity());
        existingProduct.setCategory(request.getCategory());
        existingProduct.setSku(request.getSku());

        Product updatedProduct = productRepository.save(existingProduct);
        log.info("Product updated successfully with ID: {}", updatedProduct.getId());
        
        return mapToResponse(updatedProduct);
    }

    @Override
    @Transactional
    public ProductResponse partialUpdateProduct(Long id, ProductRequest request) {
        log.info("Partial update for product with ID: {}", id);
        
        Product existingProduct = productRepository.findById(id)
            .orElseThrow(() -> new ProductNotFoundException(id));

        // Update only non-null fields
        if (request.getName() != null) {
            existingProduct.setName(request.getName());
        }
        if (request.getDescription() != null) {
            existingProduct.setDescription(request.getDescription());
        }
        if (request.getPrice() != null) {
            existingProduct.setPrice(request.getPrice());
        }
        if (request.getQuantity() != null) {
            existingProduct.setQuantity(request.getQuantity());
        }
        if (request.getCategory() != null) {
            existingProduct.setCategory(request.getCategory());
        }

        Product updatedProduct = productRepository.save(existingProduct);
        return mapToResponse(updatedProduct);
    }

    // ========== DELETE ==========
    @Override
    @Transactional
    public void deleteProduct(Long id) {
        log.info("Deleting product with ID: {}", id);
        
        if (!productRepository.existsById(id)) {
            throw new ProductNotFoundException(id);
        }
        
        productRepository.deleteById(id);
        log.info("Product deleted successfully with ID: {}", id);
    }

    @Override
    @Transactional
    public void softDeleteProduct(Long id) {
        log.info("Soft deleting product with ID: {}", id);
        
        if (!productRepository.existsById(id)) {
            throw new ProductNotFoundException(id);
        }
        
        productRepository.softDeleteById(id);
        log.info("Product soft deleted successfully with ID: {}", id);
    }

    // ========== MAPPING METHODS ==========
    private Product mapToEntity(ProductRequest request) {
        return Product.builder()
            .name(request.getName())
            .description(request.getDescription())
            .price(request.getPrice())
            .quantity(request.getQuantity())
            .category(request.getCategory())
            .sku(request.getSku())
            .active(true)
            .build();
    }

    private ProductResponse mapToResponse(Product product) {
        return ProductResponse.builder()
            .id(product.getId())
            .name(product.getName())
            .description(product.getDescription())
            .price(product.getPrice())
            .quantity(product.getQuantity())
            .category(product.getCategory())
            .sku(product.getSku())
            .active(product.getActive())
            .createdAt(product.getCreatedAt())
            .updatedAt(product.getUpdatedAt())
            .build();
    }
}
```

### 9.7 REST Controller

```java
package com.example.product.controller;

import com.example.product.dto.ProductRequest;
import com.example.product.dto.ProductResponse;
import com.example.product.service.ProductService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.web.PageableDefault;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/products")
@RequiredArgsConstructor
public class ProductController {

    private final ProductService productService;

    // ========== CREATE ==========
    @PostMapping
    public ResponseEntity<ProductResponse> createProduct(
            @Valid @RequestBody ProductRequest request) {
        ProductResponse response = productService.createProduct(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    // ========== READ ==========
    @GetMapping("/{id}")
    public ResponseEntity<ProductResponse> getProductById(@PathVariable Long id) {
        ProductResponse response = productService.getProductById(id);
        return ResponseEntity.ok(response);
    }

    @GetMapping("/sku/{sku}")
    public ResponseEntity<ProductResponse> getProductBySku(@PathVariable String sku) {
        ProductResponse response = productService.getProductBySku(sku);
        return ResponseEntity.ok(response);
    }

    @GetMapping
    public ResponseEntity<Page<ProductResponse>> getAllProducts(
            @PageableDefault(size = 10, sort = "id") Pageable pageable) {
        Page<ProductResponse> products = productService.getProducts(pageable);
        return ResponseEntity.ok(products);
    }

    @GetMapping("/list")
    public ResponseEntity<List<ProductResponse>> getAllProductsList() {
        List<ProductResponse> products = productService.getAllProducts();
        return ResponseEntity.ok(products);
    }

    @GetMapping("/category/{category}")
    public ResponseEntity<List<ProductResponse>> getProductsByCategory(
            @PathVariable String category) {
        List<ProductResponse> products = productService.getProductsByCategory(category);
        return ResponseEntity.ok(products);
    }

    // ========== UPDATE ==========
    @PutMapping("/{id}")
    public ResponseEntity<ProductResponse> updateProduct(
            @PathVariable Long id,
            @Valid @RequestBody ProductRequest request) {
        ProductResponse response = productService.updateProduct(id, request);
        return ResponseEntity.ok(response);
    }

    @PatchMapping("/{id}")
    public ResponseEntity<ProductResponse> partialUpdateProduct(
            @PathVariable Long id,
            @RequestBody ProductRequest request) {
        ProductResponse response = productService.partialUpdateProduct(id, request);
        return ResponseEntity.ok(response);
    }

    // ========== DELETE ==========
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteProduct(@PathVariable Long id) {
        productService.deleteProduct(id);
        return ResponseEntity.noContent().build();
    }

    @DeleteMapping("/{id}/soft")
    public ResponseEntity<Void> softDeleteProduct(@PathVariable Long id) {
        productService.softDeleteProduct(id);
        return ResponseEntity.noContent().build();
    }
}
```

---

### 9.8 Repository Layer Tests

```java
package com.example.product.repository;

import com.example.product.entity.Product;
import org.junit.jupiter.api.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.boot.test.autoconfigure.orm.jpa.TestEntityManager;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.test.context.ActiveProfiles;

import java.math.BigDecimal;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.*;

/**
 * Repository tests using @DataJpaTest
 * - Loads only JPA components (Entity, Repository)
 * - Uses in-memory H2 database by default
 * - Transactions are rolled back after each test
 */
@DataJpaTest
@ActiveProfiles("test")
@DisplayName("Product Repository Tests")
class ProductRepositoryTest {

    @Autowired
    private ProductRepository productRepository;

    @Autowired
    private TestEntityManager entityManager;

    private Product testProduct;

    @BeforeEach
    void setUp() {
        // Create test product using Builder pattern
        testProduct = Product.builder()
            .name("Test Product")
            .description("Test Description")
            .price(new BigDecimal("99.99"))
            .quantity(100)
            .category("Electronics")
            .sku("ELE-000001")
            .active(true)
            .build();
    }

    // ==================== CREATE TESTS ====================

    @Nested
    @DisplayName("CREATE Operations")
    class CreateTests {

        @Test
        @DisplayName("Should save product and generate ID")
        void shouldSaveProduct() {
            // When
            Product savedProduct = productRepository.save(testProduct);

            // Then
            assertThat(savedProduct).isNotNull();
            assertThat(savedProduct.getId()).isNotNull().isPositive();
            assertThat(savedProduct.getName()).isEqualTo("Test Product");
            assertThat(savedProduct.getSku()).isEqualTo("ELE-000001");
            assertThat(savedProduct.getCreatedAt()).isNotNull();
        }

        @Test
        @DisplayName("Should save product with all fields")
        void shouldSaveProductWithAllFields() {
            // Given
            Product fullProduct = Product.builder()
                .name("Full Product")
                .description("Full Description with all details")
                .price(new BigDecimal("199.99"))
                .quantity(50)
                .category("Books")
                .sku("BOK-000001")
                .active(true)
                .build();

            // When
            Product saved = productRepository.save(fullProduct);
            entityManager.flush();
            entityManager.clear();

            // Then - fetch fresh from DB
            Product fetched = productRepository.findById(saved.getId()).orElseThrow();
            assertThat(fetched.getName()).isEqualTo("Full Product");
            assertThat(fetched.getPrice()).isEqualByComparingTo(new BigDecimal("199.99"));
            assertThat(fetched.getActive()).isTrue();
        }

        @Test
        @DisplayName("Should save multiple products")
        void shouldSaveMultipleProducts() {
            // Given
            Product product2 = Product.builder()
                .name("Product 2")
                .price(new BigDecimal("49.99"))
                .quantity(200)
                .sku("ELE-000002")
                .active(true)
                .build();

            // When
            productRepository.saveAll(List.of(testProduct, product2));

            // Then
            List<Product> allProducts = productRepository.findAll();
            assertThat(allProducts).hasSize(2);
        }
    }

    // ==================== READ TESTS ====================

    @Nested
    @DisplayName("READ Operations")
    class ReadTests {

        @BeforeEach
        void setUpData() {
            entityManager.persist(testProduct);
            entityManager.flush();
        }

        @Test
        @DisplayName("Should find product by ID")
        void shouldFindById() {
            // When
            Optional<Product> found = productRepository.findById(testProduct.getId());

            // Then
            assertThat(found).isPresent();
            assertThat(found.get().getName()).isEqualTo("Test Product");
        }

        @Test
        @DisplayName("Should return empty for non-existent ID")
        void shouldReturnEmptyForNonExistentId() {
            // When
            Optional<Product> found = productRepository.findById(999L);

            // Then
            assertThat(found).isEmpty();
        }

        @Test
        @DisplayName("Should find product by SKU")
        void shouldFindBySku() {
            // When
            Optional<Product> found = productRepository.findBySku("ELE-000001");

            // Then
            assertThat(found).isPresent();
            assertThat(found.get().getName()).isEqualTo("Test Product");
        }

        @Test
        @DisplayName("Should find products by category")
        void shouldFindByCategory() {
            // Given - add more products
            Product electronics2 = Product.builder()
                .name("Electronics 2")
                .price(new BigDecimal("149.99"))
                .quantity(30)
                .category("Electronics")
                .sku("ELE-000002")
                .active(true)
                .build();
            
            Product clothing = Product.builder()
                .name("T-Shirt")
                .price(new BigDecimal("29.99"))
                .quantity(100)
                .category("Clothing")
                .sku("CLO-000001")
                .active(true)
                .build();
            
            entityManager.persist(electronics2);
            entityManager.persist(clothing);
            entityManager.flush();

            // When
            List<Product> electronics = productRepository.findByCategory("Electronics");

            // Then
            assertThat(electronics).hasSize(2);
            assertThat(electronics).extracting(Product::getCategory)
                .containsOnly("Electronics");
        }

        @Test
        @DisplayName("Should find active products only")
        void shouldFindActiveProducts() {
            // Given
            Product inactiveProduct = Product.builder()
                .name("Inactive Product")
                .price(new BigDecimal("19.99"))
                .quantity(0)
                .sku("INA-000001")
                .active(false)
                .build();
            entityManager.persist(inactiveProduct);
            entityManager.flush();

            // When
            List<Product> activeProducts = productRepository.findByActiveTrue();

            // Then
            assertThat(activeProducts).hasSize(1);
            assertThat(activeProducts.get(0).getActive()).isTrue();
        }

        @Test
        @DisplayName("Should find products by price range")
        void shouldFindByPriceRange() {
            // Given
            Product cheapProduct = Product.builder()
                .name("Cheap Product")
                .price(new BigDecimal("9.99"))
                .quantity(500)
                .sku("CHP-000001")
                .active(true)
                .build();
            
            Product expensiveProduct = Product.builder()
                .name("Expensive Product")
                .price(new BigDecimal("999.99"))
                .quantity(5)
                .sku("EXP-000001")
                .active(true)
                .build();
            
            entityManager.persist(cheapProduct);
            entityManager.persist(expensiveProduct);
            entityManager.flush();

            // When
            List<Product> midRangeProducts = productRepository.findByPriceBetween(
                new BigDecimal("50.00"), new BigDecimal("500.00"));

            // Then
            assertThat(midRangeProducts).hasSize(1);
            assertThat(midRangeProducts.get(0).getName()).isEqualTo("Test Product");
        }

        @Test
        @DisplayName("Should search products by keyword")
        void shouldSearchByKeyword() {
            // When
            List<Product> results = productRepository.searchByKeyword("Test");

            // Then
            assertThat(results).hasSize(1);
        }

        @Test
        @DisplayName("Should paginate results")
        void shouldPaginateResults() {
            // Given - add more products
            for (int i = 2; i <= 15; i++) {
                Product p = Product.builder()
                    .name("Product " + i)
                    .price(new BigDecimal("10.00"))
                    .quantity(10)
                    .sku(String.format("PRD-%06d", i))
                    .active(true)
                    .build();
                entityManager.persist(p);
            }
            entityManager.flush();

            // When
            Page<Product> page1 = productRepository.findAll(PageRequest.of(0, 5));
            Page<Product> page2 = productRepository.findAll(PageRequest.of(1, 5));

            // Then
            assertThat(page1.getContent()).hasSize(5);
            assertThat(page1.getTotalElements()).isEqualTo(15);
            assertThat(page1.getTotalPages()).isEqualTo(3);
            assertThat(page2.getContent()).hasSize(5);
        }

        @Test
        @DisplayName("Should check if SKU exists")
        void shouldCheckIfSkuExists() {
            // Then
            assertThat(productRepository.existsBySku("ELE-000001")).isTrue();
            assertThat(productRepository.existsBySku("XXX-999999")).isFalse();
        }
    }

    // ==================== UPDATE TESTS ====================

    @Nested
    @DisplayName("UPDATE Operations")
    class UpdateTests {

        @BeforeEach
        void setUpData() {
            entityManager.persist(testProduct);
            entityManager.flush();
        }

        @Test
        @DisplayName("Should update product name")
        void shouldUpdateProductName() {
            // Given
            Product product = productRepository.findById(testProduct.getId()).orElseThrow();

            // When
            product.setName("Updated Product Name");
            productRepository.save(product);
            entityManager.flush();
            entityManager.clear();

            // Then
            Product updated = productRepository.findById(testProduct.getId()).orElseThrow();
            assertThat(updated.getName()).isEqualTo("Updated Product Name");
        }

        @Test
        @DisplayName("Should update product price")
        void shouldUpdateProductPrice() {
            // Given
            Product product = productRepository.findById(testProduct.getId()).orElseThrow();

            // When
            product.setPrice(new BigDecimal("149.99"));
            productRepository.save(product);
            entityManager.flush();
            entityManager.clear();

            // Then
            Product updated = productRepository.findById(testProduct.getId()).orElseThrow();
            assertThat(updated.getPrice()).isEqualByComparingTo(new BigDecimal("149.99"));
        }

        @Test
        @DisplayName("Should update multiple fields")
        void shouldUpdateMultipleFields() {
            // Given
            Product product = productRepository.findById(testProduct.getId()).orElseThrow();

            // When
            product.setName("Multi-Update Product");
            product.setPrice(new BigDecimal("199.99"));
            product.setQuantity(50);
            product.setDescription("Updated description");
            productRepository.save(product);
            entityManager.flush();
            entityManager.clear();

            // Then
            Product updated = productRepository.findById(testProduct.getId()).orElseThrow();
            assertThat(updated.getName()).isEqualTo("Multi-Update Product");
            assertThat(updated.getPrice()).isEqualByComparingTo(new BigDecimal("199.99"));
            assertThat(updated.getQuantity()).isEqualTo(50);
            assertThat(updated.getDescription()).isEqualTo("Updated description");
            assertThat(updated.getUpdatedAt()).isNotNull();
        }

        @Test
        @DisplayName("Should soft delete product using custom query")
        void shouldSoftDeleteProduct() {
            // When
            int updatedCount = productRepository.softDeleteById(testProduct.getId());
            entityManager.flush();
            entityManager.clear();

            // Then
            assertThat(updatedCount).isEqualTo(1);
            Product updated = productRepository.findById(testProduct.getId()).orElseThrow();
            assertThat(updated.getActive()).isFalse();
        }
    }

    // ==================== DELETE TESTS ====================

    @Nested
    @DisplayName("DELETE Operations")
    class DeleteTests {

        @BeforeEach
        void setUpData() {
            entityManager.persist(testProduct);
            entityManager.flush();
        }

        @Test
        @DisplayName("Should delete product by ID")
        void shouldDeleteById() {
            // Given
            Long id = testProduct.getId();
            assertThat(productRepository.existsById(id)).isTrue();

            // When
            productRepository.deleteById(id);
            entityManager.flush();

            // Then
            assertThat(productRepository.existsById(id)).isFalse();
        }

        @Test
        @DisplayName("Should delete product entity")
        void shouldDeleteEntity() {
            // Given
            Long id = testProduct.getId();

            // When
            productRepository.delete(testProduct);
            entityManager.flush();

            // Then
            assertThat(productRepository.findById(id)).isEmpty();
        }

        @Test
        @DisplayName("Should delete all products")
        void shouldDeleteAll() {
            // Given - add more products
            Product product2 = Product.builder()
                .name("Product 2")
                .price(new BigDecimal("29.99"))
                .quantity(50)
                .sku("PRD-000002")
                .active(true)
                .build();
            entityManager.persist(product2);
            entityManager.flush();
            
            assertThat(productRepository.count()).isEqualTo(2);

            // When
            productRepository.deleteAll();
            entityManager.flush();

            // Then
            assertThat(productRepository.count()).isZero();
        }

        @Test
        @DisplayName("Should handle delete for non-existent product gracefully")
        void shouldHandleDeleteNonExistent() {
            // This test verifies behavior - deleteById doesn't throw if not found
            // But we check proper handling
            assertThat(productRepository.existsById(999L)).isFalse();
        }
    }
}
```

---

### 9.9 Service Layer Tests (Unit Tests with Mockito)

```java
package com.example.product.service;

import com.example.product.dto.ProductRequest;
import com.example.product.dto.ProductResponse;
import com.example.product.entity.Product;
import com.example.product.exception.ProductNotFoundException;
import com.example.product.repository.ProductRepository;
import org.junit.jupiter.api.*;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.*;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.*;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.*;

import static org.assertj.core.api.Assertions.*;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.BDDMockito.*;
import static org.mockito.Mockito.verify;

/**
 * Service layer unit tests using Mockito
 * - Mock repository dependencies
 * - Test business logic in isolation
 * - Verify interactions with repository
 */
@ExtendWith(MockitoExtension.class)
@DisplayName("Product Service Tests")
class ProductServiceTest {

    @Mock
    private ProductRepository productRepository;

    @InjectMocks
    private ProductServiceImpl productService;

    @Captor
    private ArgumentCaptor<Product> productCaptor;

    private Product testProduct;
    private ProductRequest testRequest;

    @BeforeEach
    void setUp() {
        testProduct = Product.builder()
            .id(1L)
            .name("Test Product")
            .description("Test Description")
            .price(new BigDecimal("99.99"))
            .quantity(100)
            .category("Electronics")
            .sku("ELE-000001")
            .active(true)
            .createdAt(LocalDateTime.now())
            .updatedAt(LocalDateTime.now())
            .build();

        testRequest = ProductRequest.builder()
            .name("Test Product")
            .description("Test Description")
            .price(new BigDecimal("99.99"))
            .quantity(100)
            .category("Electronics")
            .sku("ELE-000001")
            .build();
    }

    // ==================== CREATE TESTS ====================

    @Nested
    @DisplayName("CREATE Operations")
    class CreateTests {

        @Test
        @DisplayName("Should create product successfully")
        void shouldCreateProduct() {
            // Given (Arrange)
            given(productRepository.existsBySku(anyString())).willReturn(false);
            given(productRepository.save(any(Product.class))).willReturn(testProduct);

            // When (Act)
            ProductResponse response = productService.createProduct(testRequest);

            // Then (Assert)
            assertThat(response).isNotNull();
            assertThat(response.getId()).isEqualTo(1L);
            assertThat(response.getName()).isEqualTo("Test Product");
            assertThat(response.getSku()).isEqualTo("ELE-000001");

            // Verify interactions
            then(productRepository).should().existsBySku("ELE-000001");
            then(productRepository).should().save(productCaptor.capture());

            Product capturedProduct = productCaptor.getValue();
            assertThat(capturedProduct.getName()).isEqualTo("Test Product");
            assertThat(capturedProduct.getActive()).isTrue();
        }

        @Test
        @DisplayName("Should throw exception when SKU already exists")
        void shouldThrowExceptionWhenSkuExists() {
            // Given
            given(productRepository.existsBySku("ELE-000001")).willReturn(true);

            // When & Then
            assertThatThrownBy(() -> productService.createProduct(testRequest))
                .isInstanceOf(IllegalArgumentException.class)
                .hasMessageContaining("already exists");

            // Verify save was never called
            then(productRepository).should(never()).save(any());
        }

        @Test
        @DisplayName("Should map all fields correctly when creating product")
        void shouldMapAllFieldsCorrectly() {
            // Given
            given(productRepository.existsBySku(anyString())).willReturn(false);
            given(productRepository.save(any(Product.class))).willAnswer(invocation -> {
                Product p = invocation.getArgument(0);
                p.setId(1L);
                p.setCreatedAt(LocalDateTime.now());
                return p;
            });

            // When
            productService.createProduct(testRequest);

            // Then - verify entity mapping
            then(productRepository).should().save(productCaptor.capture());
            Product captured = productCaptor.getValue();
            
            assertThat(captured.getName()).isEqualTo(testRequest.getName());
            assertThat(captured.getDescription()).isEqualTo(testRequest.getDescription());
            assertThat(captured.getPrice()).isEqualTo(testRequest.getPrice());
            assertThat(captured.getQuantity()).isEqualTo(testRequest.getQuantity());
            assertThat(captured.getCategory()).isEqualTo(testRequest.getCategory());
            assertThat(captured.getSku()).isEqualTo(testRequest.getSku());
        }
    }

    // ==================== READ TESTS ====================

    @Nested
    @DisplayName("READ Operations")
    class ReadTests {

        @Test
        @DisplayName("Should find product by ID")
        void shouldFindProductById() {
            // Given
            given(productRepository.findById(1L)).willReturn(Optional.of(testProduct));

            // When
            ProductResponse response = productService.getProductById(1L);

            // Then
            assertThat(response).isNotNull();
            assertThat(response.getId()).isEqualTo(1L);
            assertThat(response.getName()).isEqualTo("Test Product");

            then(productRepository).should().findById(1L);
        }

        @Test
        @DisplayName("Should throw ProductNotFoundException when ID not found")
        void shouldThrowNotFoundExceptionForInvalidId() {
            // Given
            given(productRepository.findById(999L)).willReturn(Optional.empty());

            // When & Then
            assertThatThrownBy(() -> productService.getProductById(999L))
                .isInstanceOf(ProductNotFoundException.class)
                .hasMessageContaining("999");
        }

        @Test
        @DisplayName("Should find product by SKU")
        void shouldFindProductBySku() {
            // Given
            given(productRepository.findBySku("ELE-000001")).willReturn(Optional.of(testProduct));

            // When
            ProductResponse response = productService.getProductBySku("ELE-000001");

            // Then
            assertThat(response).isNotNull();
            assertThat(response.getSku()).isEqualTo("ELE-000001");
        }

        @Test
        @DisplayName("Should throw ProductNotFoundException when SKU not found")
        void shouldThrowNotFoundExceptionForInvalidSku() {
            // Given
            given(productRepository.findBySku("XXX-999999")).willReturn(Optional.empty());

            // When & Then
            assertThatThrownBy(() -> productService.getProductBySku("XXX-999999"))
                .isInstanceOf(ProductNotFoundException.class)
                .hasMessageContaining("XXX-999999");
        }

        @Test
        @DisplayName("Should return all products")
        void shouldReturnAllProducts() {
            // Given
            Product product2 = Product.builder()
                .id(2L)
                .name("Product 2")
                .price(new BigDecimal("49.99"))
                .quantity(50)
                .sku("ELE-000002")
                .active(true)
                .build();
            
            given(productRepository.findAll()).willReturn(List.of(testProduct, product2));

            // When
            List<ProductResponse> products = productService.getAllProducts();

            // Then
            assertThat(products).hasSize(2);
            assertThat(products).extracting(ProductResponse::getId).containsExactly(1L, 2L);
        }

        @Test
        @DisplayName("Should return empty list when no products exist")
        void shouldReturnEmptyList() {
            // Given
            given(productRepository.findAll()).willReturn(Collections.emptyList());

            // When
            List<ProductResponse> products = productService.getAllProducts();

            // Then
            assertThat(products).isEmpty();
        }

        @Test
        @DisplayName("Should find products by category")
        void shouldFindProductsByCategory() {
            // Given
            given(productRepository.findByCategory("Electronics"))
                .willReturn(List.of(testProduct));

            // When
            List<ProductResponse> products = productService.getProductsByCategory("Electronics");

            // Then
            assertThat(products).hasSize(1);
            assertThat(products.get(0).getCategory()).isEqualTo("Electronics");
        }

        @Test
        @DisplayName("Should return paginated products")
        void shouldReturnPaginatedProducts() {
            // Given
            Pageable pageable = PageRequest.of(0, 10);
            Page<Product> productPage = new PageImpl<>(
                List.of(testProduct),
                pageable,
                1
            );
            given(productRepository.findAll(pageable)).willReturn(productPage);

            // When
            Page<ProductResponse> responsePage = productService.getProducts(pageable);

            // Then
            assertThat(responsePage.getContent()).hasSize(1);
            assertThat(responsePage.getTotalElements()).isEqualTo(1);
            assertThat(responsePage.getNumber()).isZero();
        }
    }

    // ==================== UPDATE TESTS ====================

    @Nested
    @DisplayName("UPDATE Operations")
    class UpdateTests {

        @Test
        @DisplayName("Should update product successfully")
        void shouldUpdateProduct() {
            // Given
            ProductRequest updateRequest = ProductRequest.builder()
                .name("Updated Product")
                .description("Updated Description")
                .price(new BigDecimal("149.99"))
                .quantity(200)
                .category("Electronics")
                .sku("ELE-000001")  // Same SKU
                .build();

            given(productRepository.findById(1L)).willReturn(Optional.of(testProduct));
            given(productRepository.save(any(Product.class))).willAnswer(invocation -> invocation.getArgument(0));

            // When
            ProductResponse response = productService.updateProduct(1L, updateRequest);

            // Then
            assertThat(response.getName()).isEqualTo("Updated Product");
            assertThat(response.getPrice()).isEqualByComparingTo(new BigDecimal("149.99"));
            assertThat(response.getQuantity()).isEqualTo(200);

            then(productRepository).should().save(productCaptor.capture());
            Product saved = productCaptor.getValue();
            assertThat(saved.getName()).isEqualTo("Updated Product");
        }

        @Test
        @DisplayName("Should throw exception when updating non-existent product")
        void shouldThrowExceptionWhenUpdatingNonExistent() {
            // Given
            given(productRepository.findById(999L)).willReturn(Optional.empty());

            // When & Then
            assertThatThrownBy(() -> productService.updateProduct(999L, testRequest))
                .isInstanceOf(ProductNotFoundException.class);

            then(productRepository).should(never()).save(any());
        }

        @Test
        @DisplayName("Should throw exception when updating with duplicate SKU")
        void shouldThrowExceptionForDuplicateSku() {
            // Given
            ProductRequest updateRequest = ProductRequest.builder()
                .name("Updated Product")
                .price(new BigDecimal("99.99"))
                .quantity(100)
                .sku("ELE-000002")  // Different SKU that exists
                .build();

            given(productRepository.findById(1L)).willReturn(Optional.of(testProduct));
            given(productRepository.existsBySku("ELE-000002")).willReturn(true);

            // When & Then
            assertThatThrownBy(() -> productService.updateProduct(1L, updateRequest))
                .isInstanceOf(IllegalArgumentException.class)
                .hasMessageContaining("already exists");
        }

        @Test
        @DisplayName("Should allow keeping same SKU during update")
        void shouldAllowSameSkuDuringUpdate() {
            // Given - request with same SKU
            given(productRepository.findById(1L)).willReturn(Optional.of(testProduct));
            given(productRepository.save(any(Product.class))).willReturn(testProduct);

            // When - existsBySku should NOT be called for same SKU
            ProductResponse response = productService.updateProduct(1L, testRequest);

            // Then
            assertThat(response).isNotNull();
            then(productRepository).should(never()).existsBySku(anyString());
        }

        @Test
        @DisplayName("Should partially update product - only non-null fields")
        void shouldPartiallyUpdateProduct() {
            // Given
            ProductRequest partialRequest = ProductRequest.builder()
                .name("Partially Updated")
                .price(new BigDecimal("199.99"))
                // Other fields are null
                .build();

            given(productRepository.findById(1L)).willReturn(Optional.of(testProduct));
            given(productRepository.save(any(Product.class))).willAnswer(invocation -> invocation.getArgument(0));

            // When
            ProductResponse response = productService.partialUpdateProduct(1L, partialRequest);

            // Then
            then(productRepository).should().save(productCaptor.capture());
            Product saved = productCaptor.getValue();
            
            // Updated fields
            assertThat(saved.getName()).isEqualTo("Partially Updated");
            assertThat(saved.getPrice()).isEqualByComparingTo(new BigDecimal("199.99"));
            
            // Unchanged fields (original values)
            assertThat(saved.getDescription()).isEqualTo("Test Description");
            assertThat(saved.getQuantity()).isEqualTo(100);
            assertThat(saved.getCategory()).isEqualTo("Electronics");
        }
    }

    // ==================== DELETE TESTS ====================

    @Nested
    @DisplayName("DELETE Operations")
    class DeleteTests {

        @Test
        @DisplayName("Should delete product successfully")
        void shouldDeleteProduct() {
            // Given
            given(productRepository.existsById(1L)).willReturn(true);
            willDoNothing().given(productRepository).deleteById(1L);

            // When
            productService.deleteProduct(1L);

            // Then
            then(productRepository).should().existsById(1L);
            then(productRepository).should().deleteById(1L);
        }

        @Test
        @DisplayName("Should throw exception when deleting non-existent product")
        void shouldThrowExceptionWhenDeletingNonExistent() {
            // Given
            given(productRepository.existsById(999L)).willReturn(false);

            // When & Then
            assertThatThrownBy(() -> productService.deleteProduct(999L))
                .isInstanceOf(ProductNotFoundException.class);

            then(productRepository).should(never()).deleteById(anyLong());
        }

        @Test
        @DisplayName("Should soft delete product successfully")
        void shouldSoftDeleteProduct() {
            // Given
            given(productRepository.existsById(1L)).willReturn(true);
            given(productRepository.softDeleteById(1L)).willReturn(1);

            // When
            productService.softDeleteProduct(1L);

            // Then
            then(productRepository).should().softDeleteById(1L);
            then(productRepository).should(never()).deleteById(anyLong());
        }

        @Test
        @DisplayName("Should throw exception when soft deleting non-existent product")
        void shouldThrowExceptionWhenSoftDeletingNonExistent() {
            // Given
            given(productRepository.existsById(999L)).willReturn(false);

            // When & Then
            assertThatThrownBy(() -> productService.softDeleteProduct(999L))
                .isInstanceOf(ProductNotFoundException.class);

            then(productRepository).should(never()).softDeleteById(anyLong());
        }
    }

    // ==================== VERIFICATION TESTS ====================

    @Nested
    @DisplayName("Verification Patterns")
    class VerificationTests {

        @Test
        @DisplayName("Should verify exact number of repository calls")
        void shouldVerifyExactCalls() {
            // Given
            given(productRepository.findById(1L)).willReturn(Optional.of(testProduct));
            given(productRepository.findById(2L)).willReturn(Optional.empty());

            // When
            productService.getProductById(1L);
            try {
                productService.getProductById(2L);
            } catch (ProductNotFoundException ignored) {}

            // Then - verify findById was called exactly 2 times
            then(productRepository).should(times(2)).findById(anyLong());
        }

        @Test
        @DisplayName("Should verify no more interactions after operation")
        void shouldVerifyNoMoreInteractions() {
            // Given
            given(productRepository.existsById(1L)).willReturn(true);
            willDoNothing().given(productRepository).deleteById(1L);

            // When
            productService.deleteProduct(1L);

            // Then
            then(productRepository).should().existsById(1L);
            then(productRepository).should().deleteById(1L);
            then(productRepository).shouldHaveNoMoreInteractions();
        }

        @Test
        @DisplayName("Should verify order of method calls")
        void shouldVerifyCallOrder() {
            // Given
            given(productRepository.existsBySku(anyString())).willReturn(false);
            given(productRepository.save(any(Product.class))).willReturn(testProduct);

            // When
            productService.createProduct(testRequest);

            // Then - verify order: first check SKU exists, then save
            InOrder inOrder = inOrder(productRepository);
            inOrder.verify(productRepository).existsBySku("ELE-000001");
            inOrder.verify(productRepository).save(any(Product.class));
        }
    }
}
```

---

### 9.10 Controller Layer Tests (Integration with MockMvc)

```java
package com.example.product.controller;

import com.example.product.dto.ProductRequest;
import com.example.product.dto.ProductResponse;
import com.example.product.exception.GlobalExceptionHandler;
import com.example.product.exception.ProductNotFoundException;
import com.example.product.service.ProductService;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.*;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.data.domain.*;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.ResultActions;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.*;

import static org.hamcrest.Matchers.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.BDDMockito.*;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

/**
 * Controller layer tests using @WebMvcTest
 * - Loads only web layer (Controller, ControllerAdvice)
 * - Service is mocked using @MockBean
 * - Tests HTTP requests/responses, validation, status codes
 */
@WebMvcTest(ProductController.class)
@DisplayName("Product Controller Tests")
class ProductControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @MockBean
    private ProductService productService;

    private ProductRequest validRequest;
    private ProductResponse sampleResponse;

    @BeforeEach
    void setUp() {
        validRequest = ProductRequest.builder()
            .name("Test Product")
            .description("Test Description")
            .price(new BigDecimal("99.99"))
            .quantity(100)
            .category("Electronics")
            .sku("ELE-000001")
            .build();

        sampleResponse = ProductResponse.builder()
            .id(1L)
            .name("Test Product")
            .description("Test Description")
            .price(new BigDecimal("99.99"))
            .quantity(100)
            .category("Electronics")
            .sku("ELE-000001")
            .active(true)
            .createdAt(LocalDateTime.now())
            .updatedAt(LocalDateTime.now())
            .build();
    }

    // ==================== CREATE TESTS ====================

    @Nested
    @DisplayName("POST /api/v1/products")
    class CreateProductTests {

        @Test
        @DisplayName("Should create product and return 201 CREATED")
        void shouldCreateProductSuccessfully() throws Exception {
            // Given
            given(productService.createProduct(any(ProductRequest.class)))
                .willReturn(sampleResponse);

            // When
            ResultActions result = mockMvc.perform(post("/api/v1/products")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsString(validRequest)));

            // Then
            result.andDo(print())
                .andExpect(status().isCreated())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$.id", is(1)))
                .andExpect(jsonPath("$.name", is("Test Product")))
                .andExpect(jsonPath("$.sku", is("ELE-000001")))
                .andExpect(jsonPath("$.active", is(true)));

            then(productService).should().createProduct(any(ProductRequest.class));
        }

        @Test
        @DisplayName("Should return 400 BAD REQUEST for missing name")
        void shouldReturn400ForMissingName() throws Exception {
            // Given
            ProductRequest invalidRequest = ProductRequest.builder()
                .name(null)  // Missing required field
                .price(new BigDecimal("99.99"))
                .quantity(100)
                .sku("ELE-000001")
                .build();

            // When & Then
            mockMvc.perform(post("/api/v1/products")
                    .contentType(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(invalidRequest)))
                .andDo(print())
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.errors.name").exists());

            then(productService).should(never()).createProduct(any());
        }

        @Test
        @DisplayName("Should return 400 BAD REQUEST for invalid price")
        void shouldReturn400ForInvalidPrice() throws Exception {
            // Given
            ProductRequest invalidRequest = ProductRequest.builder()
                .name("Test Product")
                .price(new BigDecimal("-10.00"))  // Invalid: negative price
                .quantity(100)
                .sku("ELE-000001")
                .build();

            // When & Then
            mockMvc.perform(post("/api/v1/products")
                    .contentType(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(invalidRequest)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.errors.price").exists());
        }

        @Test
        @DisplayName("Should return 400 BAD REQUEST for invalid SKU format")
        void shouldReturn400ForInvalidSku() throws Exception {
            // Given
            ProductRequest invalidRequest = ProductRequest.builder()
                .name("Test Product")
                .price(new BigDecimal("99.99"))
                .quantity(100)
                .sku("INVALID-SKU")  // Invalid format
                .build();

            // When & Then
            mockMvc.perform(post("/api/v1/products")
                    .contentType(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(invalidRequest)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.errors.sku").exists());
        }

        @Test
        @DisplayName("Should return 400 BAD REQUEST for negative quantity")
        void shouldReturn400ForNegativeQuantity() throws Exception {
            // Given
            ProductRequest invalidRequest = ProductRequest.builder()
                .name("Test Product")
                .price(new BigDecimal("99.99"))
                .quantity(-5)  // Invalid: negative
                .sku("ELE-000001")
                .build();

            // When & Then
            mockMvc.perform(post("/api/v1/products")
                    .contentType(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(invalidRequest)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.errors.quantity").exists());
        }

        @Test
        @DisplayName("Should return 400 for multiple validation errors")
        void shouldReturn400ForMultipleErrors() throws Exception {
            // Given - empty request with multiple missing fields
            ProductRequest invalidRequest = ProductRequest.builder().build();

            // When & Then
            mockMvc.perform(post("/api/v1/products")
                    .contentType(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(invalidRequest)))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.errors").isMap())
                .andExpect(jsonPath("$.errors", aMapWithSize(greaterThanOrEqualTo(3))));
        }
    }

    // ==================== READ TESTS ====================

    @Nested
    @DisplayName("GET /api/v1/products")
    class GetProductTests {

        @Test
        @DisplayName("Should get product by ID and return 200 OK")
        void shouldGetProductById() throws Exception {
            // Given
            given(productService.getProductById(1L)).willReturn(sampleResponse);

            // When & Then
            mockMvc.perform(get("/api/v1/products/{id}", 1L))
                .andDo(print())
                .andExpect(status().isOk())
                .andExpect(content().contentType(MediaType.APPLICATION_JSON))
                .andExpect(jsonPath("$.id", is(1)))
                .andExpect(jsonPath("$.name", is("Test Product")))
                .andExpect(jsonPath("$.price", is(99.99)));
        }

        @Test
        @DisplayName("Should return 404 NOT FOUND for non-existent ID")
        void shouldReturn404ForNonExistentId() throws Exception {
            // Given
            given(productService.getProductById(999L))
                .willThrow(new ProductNotFoundException(999L));

            // When & Then
            mockMvc.perform(get("/api/v1/products/{id}", 999L))
                .andExpect(status().isNotFound())
                .andExpect(jsonPath("$.message", containsString("999")));
        }

        @Test
        @DisplayName("Should get product by SKU")
        void shouldGetProductBySku() throws Exception {
            // Given
            given(productService.getProductBySku("ELE-000001")).willReturn(sampleResponse);

            // When & Then
            mockMvc.perform(get("/api/v1/products/sku/{sku}", "ELE-000001"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.sku", is("ELE-000001")));
        }

        @Test
        @DisplayName("Should get all products as list")
        void shouldGetAllProductsList() throws Exception {
            // Given
            ProductResponse product2 = ProductResponse.builder()
                .id(2L)
                .name("Product 2")
                .price(new BigDecimal("49.99"))
                .quantity(50)
                .sku("ELE-000002")
                .active(true)
                .build();

            given(productService.getAllProducts())
                .willReturn(List.of(sampleResponse, product2));

            // When & Then
            mockMvc.perform(get("/api/v1/products/list"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(2)))
                .andExpect(jsonPath("$[0].id", is(1)))
                .andExpect(jsonPath("$[1].id", is(2)));
        }

        @Test
        @DisplayName("Should return empty list when no products")
        void shouldReturnEmptyList() throws Exception {
            // Given
            given(productService.getAllProducts()).willReturn(Collections.emptyList());

            // When & Then
            mockMvc.perform(get("/api/v1/products/list"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(0)));
        }

        @Test
        @DisplayName("Should get paginated products")
        void shouldGetPaginatedProducts() throws Exception {
            // Given
            Page<ProductResponse> productPage = new PageImpl<>(
                List.of(sampleResponse),
                PageRequest.of(0, 10),
                1
            );
            given(productService.getProducts(any(Pageable.class))).willReturn(productPage);

            // When & Then
            mockMvc.perform(get("/api/v1/products")
                    .param("page", "0")
                    .param("size", "10"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(1)))
                .andExpect(jsonPath("$.totalElements", is(1)))
                .andExpect(jsonPath("$.totalPages", is(1)));
        }

        @Test
        @DisplayName("Should get products by category")
        void shouldGetProductsByCategory() throws Exception {
            // Given
            given(productService.getProductsByCategory("Electronics"))
                .willReturn(List.of(sampleResponse));

            // When & Then
            mockMvc.perform(get("/api/v1/products/category/{category}", "Electronics"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$", hasSize(1)))
                .andExpect(jsonPath("$[0].category", is("Electronics")));
        }
    }

    // ==================== UPDATE TESTS ====================

    @Nested
    @DisplayName("PUT/PATCH /api/v1/products/{id}")
    class UpdateProductTests {

        @Test
        @DisplayName("Should update product and return 200 OK")
        void shouldUpdateProductSuccessfully() throws Exception {
            // Given
            ProductRequest updateRequest = ProductRequest.builder()
                .name("Updated Product")
                .description("Updated Description")
                .price(new BigDecimal("149.99"))
                .quantity(200)
                .category("Electronics")
                .sku("ELE-000001")
                .build();

            ProductResponse updatedResponse = ProductResponse.builder()
                .id(1L)
                .name("Updated Product")
                .description("Updated Description")
                .price(new BigDecimal("149.99"))
                .quantity(200)
                .category("Electronics")
                .sku("ELE-000001")
                .active(true)
                .build();

            given(productService.updateProduct(eq(1L), any(ProductRequest.class)))
                .willReturn(updatedResponse);

            // When & Then
            mockMvc.perform(put("/api/v1/products/{id}", 1L)
                    .contentType(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(updateRequest)))
                .andDo(print())
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.name", is("Updated Product")))
                .andExpect(jsonPath("$.price", is(149.99)))
                .andExpect(jsonPath("$.quantity", is(200)));
        }

        @Test
        @DisplayName("Should return 404 when updating non-existent product")
        void shouldReturn404WhenUpdatingNonExistent() throws Exception {
            // Given
            given(productService.updateProduct(eq(999L), any(ProductRequest.class)))
                .willThrow(new ProductNotFoundException(999L));

            // When & Then
            mockMvc.perform(put("/api/v1/products/{id}", 999L)
                    .contentType(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(validRequest)))
                .andExpect(status().isNotFound());
        }

        @Test
        @DisplayName("Should return 400 for invalid update data")
        void shouldReturn400ForInvalidUpdateData() throws Exception {
            // Given
            ProductRequest invalidRequest = ProductRequest.builder()
                .name("")  // Invalid: blank name
                .price(new BigDecimal("99.99"))
                .quantity(100)
                .sku("ELE-000001")
                .build();

            // When & Then
            mockMvc.perform(put("/api/v1/products/{id}", 1L)
                    .contentType(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(invalidRequest)))
                .andExpect(status().isBadRequest());
        }

        @Test
        @DisplayName("Should partially update product using PATCH")
        void shouldPartiallyUpdateProduct() throws Exception {
            // Given
            ProductRequest partialRequest = ProductRequest.builder()
                .name("Partially Updated")
                .build();

            ProductResponse partiallyUpdated = ProductResponse.builder()
                .id(1L)
                .name("Partially Updated")
                .description("Test Description")  // Unchanged
                .price(new BigDecimal("99.99"))    // Unchanged
                .quantity(100)                     // Unchanged
                .category("Electronics")          // Unchanged
                .sku("ELE-000001")
                .active(true)
                .build();

            given(productService.partialUpdateProduct(eq(1L), any(ProductRequest.class)))
                .willReturn(partiallyUpdated);

            // When & Then
            mockMvc.perform(patch("/api/v1/products/{id}", 1L)
                    .contentType(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(partialRequest)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.name", is("Partially Updated")))
                .andExpect(jsonPath("$.price", is(99.99)));  // Unchanged
        }
    }

    // ==================== DELETE TESTS ====================

    @Nested
    @DisplayName("DELETE /api/v1/products/{id}")
    class DeleteProductTests {

        @Test
        @DisplayName("Should delete product and return 204 NO CONTENT")
        void shouldDeleteProductSuccessfully() throws Exception {
            // Given
            willDoNothing().given(productService).deleteProduct(1L);

            // When & Then
            mockMvc.perform(delete("/api/v1/products/{id}", 1L))
                .andExpect(status().isNoContent())
                .andExpect(content().string(""));  // No body

            then(productService).should().deleteProduct(1L);
        }

        @Test
        @DisplayName("Should return 404 when deleting non-existent product")
        void shouldReturn404WhenDeletingNonExistent() throws Exception {
            // Given
            willThrow(new ProductNotFoundException(999L))
                .given(productService).deleteProduct(999L);

            // When & Then
            mockMvc.perform(delete("/api/v1/products/{id}", 999L))
                .andExpect(status().isNotFound());
        }

        @Test
        @DisplayName("Should soft delete product")
        void shouldSoftDeleteProduct() throws Exception {
            // Given
            willDoNothing().given(productService).softDeleteProduct(1L);

            // When & Then
            mockMvc.perform(delete("/api/v1/products/{id}/soft", 1L))
                .andExpect(status().isNoContent());

            then(productService).should().softDeleteProduct(1L);
        }
    }

    // ==================== CONTENT TYPE & HEADERS TESTS ====================

    @Nested
    @DisplayName("Content Type & Headers")
    class ContentTypeTests {

        @Test
        @DisplayName("Should return 415 for unsupported media type")
        void shouldReturn415ForUnsupportedMediaType() throws Exception {
            mockMvc.perform(post("/api/v1/products")
                    .contentType(MediaType.TEXT_PLAIN)
                    .content("invalid content"))
                .andExpect(status().isUnsupportedMediaType());
        }

        @Test
        @DisplayName("Should accept application/json content type")
        void shouldAcceptJsonContentType() throws Exception {
            // Given
            given(productService.createProduct(any(ProductRequest.class)))
                .willReturn(sampleResponse);

            // When & Then
            mockMvc.perform(post("/api/v1/products")
                    .contentType(MediaType.APPLICATION_JSON)
                    .accept(MediaType.APPLICATION_JSON)
                    .content(objectMapper.writeValueAsString(validRequest)))
                .andExpect(status().isCreated())
                .andExpect(header().string("Content-Type", containsString("application/json")));
        }
    }
}
```

---

### 9.11 Test Configuration & Utilities

#### Test Properties (application-test.yml)

```yaml
# src/test/resources/application-test.yml
spring:
  datasource:
    url: jdbc:h2:mem:testdb;DB_CLOSE_DELAY=-1;DB_CLOSE_ON_EXIT=FALSE
    driver-class-name: org.h2.Driver
    username: sa
    password:
  
  jpa:
    hibernate:
      ddl-auto: create-drop
    show-sql: true
    properties:
      hibernate:
        format_sql: true
  
  h2:
    console:
      enabled: false

logging:
  level:
    org.springframework.test: INFO
    com.example.product: DEBUG
    org.hibernate.SQL: DEBUG
```

#### Test Data Builder (Optional Clean Pattern)

```java
package com.example.product.util;

import com.example.product.dto.ProductRequest;
import com.example.product.dto.ProductResponse;
import com.example.product.entity.Product;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.concurrent.atomic.AtomicLong;

/**
 * Test data builder for creating test fixtures
 * Provides sensible defaults while allowing customization
 */
public class TestDataBuilder {

    private static final AtomicLong ID_COUNTER = new AtomicLong(1);
    private static final AtomicLong SKU_COUNTER = new AtomicLong(1);

    // ==================== Product Entity Builders ====================

    public static Product.ProductBuilder defaultProduct() {
        long skuNum = SKU_COUNTER.getAndIncrement();
        return Product.builder()
            .name("Test Product")
            .description("Test Description")
            .price(new BigDecimal("99.99"))
            .quantity(100)
            .category("Electronics")
            .sku(String.format("TST-%06d", skuNum))
            .active(true)
            .createdAt(LocalDateTime.now())
            .updatedAt(LocalDateTime.now());
    }

    public static Product.ProductBuilder productWithId() {
        return defaultProduct().id(ID_COUNTER.getAndIncrement());
    }

    public static Product.ProductBuilder electronics() {
        return defaultProduct().category("Electronics");
    }

    public static Product.ProductBuilder clothing() {
        return defaultProduct()
            .category("Clothing")
            .price(new BigDecimal("29.99"));
    }

    public static Product.ProductBuilder inactiveProduct() {
        return defaultProduct().active(false);
    }

    // ==================== ProductRequest Builders ====================

    public static ProductRequest.ProductRequestBuilder defaultRequest() {
        long skuNum = SKU_COUNTER.getAndIncrement();
        return ProductRequest.builder()
            .name("Test Product")
            .description("Test Description")
            .price(new BigDecimal("99.99"))
            .quantity(100)
            .category("Electronics")
            .sku(String.format("TST-%06d", skuNum));
    }

    public static ProductRequest.ProductRequestBuilder invalidRequest() {
        return ProductRequest.builder();
    }

    // ==================== ProductResponse Builders ====================

    public static ProductResponse.ProductResponseBuilder defaultResponse() {
        long id = ID_COUNTER.getAndIncrement();
        long skuNum = SKU_COUNTER.getAndIncrement();
        return ProductResponse.builder()
            .id(id)
            .name("Test Product")
            .description("Test Description")
            .price(new BigDecimal("99.99"))
            .quantity(100)
            .category("Electronics")
            .sku(String.format("TST-%06d", skuNum))
            .active(true)
            .createdAt(LocalDateTime.now())
            .updatedAt(LocalDateTime.now());
    }

    // ==================== Reset Methods ====================

    public static void resetCounters() {
        ID_COUNTER.set(1);
        SKU_COUNTER.set(1);
    }
}
```

#### Using Test Builder in Tests

```java
import static com.example.product.util.TestDataBuilder.*;

@BeforeEach
void setUp() {
    resetCounters();
}

@Test
void shouldCreateProduct() {
    // Given - clean, readable test data
    Product product = electronics().name("Laptop").price(new BigDecimal("999.99")).build();
    ProductRequest request = defaultRequest().name("Laptop").build();
    
    // ...
}
```

---

### 9.12 Summary - Testing Cheat Sheet

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        CRUD Testing Layers Summary                               │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  LAYER           ANNOTATION          WHAT'S TESTED           DEPENDENCIES        │
│  ─────────────────────────────────────────────────────────────────────────────  │
│  Repository      @DataJpaTest        JPA queries, CRUD       H2/Testcontainers  │
│                                      Entity mappings         Auto-rollback      │
│                                                                                  │
│  Service         @ExtendWith         Business logic          @Mock Repository   │
│                  (MockitoExtension)  Validation rules        @InjectMocks       │
│                                      Exception handling                          │
│                                                                                  │
│  Controller      @WebMvcTest         HTTP status codes       @MockBean Service  │
│                                      Request validation      MockMvc            │
│                                      JSON serialization                          │
│                                                                                  │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  TEST STRUCTURE (AAA Pattern):                                                   │
│  ┌──────────────────────────────────────────────────────────────────────────┐  │
│  │  @Test                                                                    │  │
│  │  void shouldDoSomething() {                                              │  │
│  │      // Given (Arrange) - Set up test data & mocks                       │  │
│  │      given(repository.findById(1L)).willReturn(Optional.of(entity));    │  │
│  │                                                                          │  │
│  │      // When (Act) - Execute the method under test                       │  │
│  │      Result result = service.doSomething(1L);                           │  │
│  │                                                                          │  │
│  │      // Then (Assert) - Verify results & interactions                    │  │
│  │      assertThat(result).isNotNull();                                    │  │
│  │      then(repository).should().findById(1L);                            │  │
│  │  }                                                                       │  │
│  └──────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  KEY ASSERTIONS:                                                                 │
│  ─────────────────────────────────────────────────────────────────────────────  │
│  assertThat(result).isNotNull();                                                │
│  assertThat(result.getName()).isEqualTo("expected");                           │
│  assertThat(list).hasSize(3);                                                   │
│  assertThatThrownBy(() -> service.method()).isInstanceOf(Exception.class);     │
│                                                                                  │
│  MOCK SETUP:                                                                     │
│  ─────────────────────────────────────────────────────────────────────────────  │
│  given(mock.method()).willReturn(value);      // Return value                   │
│  given(mock.method()).willThrow(exception);   // Throw exception                │
│  willDoNothing().given(mock).voidMethod();    // Void methods                   │
│                                                                                  │
│  MOCK VERIFICATION:                                                              │
│  ─────────────────────────────────────────────────────────────────────────────  │
│  then(mock).should().method();                // Called once                    │
│  then(mock).should(times(2)).method();        // Called twice                   │
│  then(mock).should(never()).method();         // Never called                   │
│                                                                                  │
│  HTTP STATUS CODES (Controller Tests):                                          │
│  ─────────────────────────────────────────────────────────────────────────────  │
│  200 OK          - GET successful                                               │
│  201 CREATED     - POST successful                                              │
│  204 NO CONTENT  - DELETE successful                                            │
│  400 BAD REQUEST - Validation failed                                            │
│  404 NOT FOUND   - Resource not found                                           │
│  415 UNSUPPORTED - Wrong Content-Type                                           │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Quick Reference

### Common Annotations

| Annotation | Purpose |
|------------|---------|
| `@Test` | Marks test method |
| `@BeforeEach` / `@AfterEach` | Before/after each test |
| `@BeforeAll` / `@AfterAll` | Before/after all tests |
| `@DisplayName` | Custom test name |
| `@Disabled` | Skip test |
| `@Mock` | Create mock |
| `@InjectMocks` | Inject mocks into class |
| `@Spy` | Partial mock |
| `@MockBean` | Replace Spring bean with mock |
| `@SpyBean` | Wrap Spring bean with spy |

### Common Matchers

```java
// ArgumentMatchers
any(), anyString(), anyLong(), anyList()
eq(value), isNull(), isNotNull()
argThat(predicate)

// Hamcrest (for assertions)
is(), equalTo(), hasSize(), containsString()
hasItem(), hasItems(), containsInAnyOrder()
greaterThan(), lessThan(), closeTo()
```

### Verification Modes

```java
times(n)        // Exactly n times
never()         // Zero times
atLeast(n)      // At least n times
atMost(n)       // At most n times
atLeastOnce()   // At least once
```

---

*This completes the JUnit 5 & Mockito Testing Guide.*

*Last Updated: February 2026*
