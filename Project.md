Here’s a strong, ready-to-use prompt you can give to an AI (or use as a project brief) to build your stock price prediction system:

---

**Prompt:**

> Design and develop a comprehensive, data-driven stock price prediction system that identifies stocks likely to experience high price movements (volatility) within intraday or short-term periods (1–5 days).
>
> The system should integrate multiple data sources, including:
>
> * Historical stock price data (OHLCV)
> * Financial reports (earnings, balance sheets, key ratios)
> * Real-time and historical news data
> * Market sentiment (news sentiment, social media if possible)
> * Macroeconomic indicators (interest rates, inflation, global indices)
>
> Key requirements:
>
> 1. **Data Processing & Feature Engineering**
>
>    * Clean and preprocess structured and unstructured data
>    * Extract features such as technical indicators (RSI, MACD, moving averages, volume spikes)
>    * Perform NLP-based sentiment analysis on news and reports
>    * Detect anomalies or unusual patterns in volume and price
> 2. **Prediction Goals**
>
>    * Predict probability of high movement (e.g., >2–5%) in the next trading session or next few days
>    * Classify stocks into categories: High Movement, Moderate Movement, Low Movement
>    * Optionally predict direction (bullish/bearish)
> 3. **Modeling Approach**
>
>    * Use a combination of:
>
>      * Time-series models (LSTM, GRU, ARIMA)
>      * Machine learning models (XGBoost, Random Forest)
>      * NLP models (transformers for sentiment analysis)
>    * Ensemble multiple models for improved accuracy
> 4. **Event & News Impact Analysis**
>
>    * Detect key events (earnings announcements, mergers, policy changes)
>    * Quantify expected impact of news sentiment on price movement
> 5. **Output & Decision Support**
>
>    * Rank stocks based on predicted movement probability
>    * Provide confidence scores and key reasons (explainability)
>    * Generate actionable insights (e.g., “High volatility expected due to earnings + positive sentiment + volume spike”)
> 6. **Evaluation Metrics**
>
>    * Accuracy, precision/recall for classification
>    * Backtesting on historical data
>    * Sharpe ratio or return-based evaluation for strategy validation
> 7. **System Design**
>
>    * Real-time data ingestion pipeline
>    * Scalable architecture (APIs, dashboards)
>    * Visualization of predictions and trends
> 8. **Constraints & Considerations**
>
>    * Avoid data leakage and ensure proper train/test splits
>    * Handle noisy and delayed news data
>    * Ensure model interpretability for decision-making
>
> Deliver:
>
> * End-to-end architecture
> * Model pipeline
> * Example implementation (Python preferred)
> * Suggestions for deployment (cloud/local)
>
> The system should prioritize **data-driven decision making**, robustness, and practical usability for traders or analysts.

---

If you want, I can refine this into:

* a **step-by-step implementation plan**,
* a **resume/project description**, or
* a **production-grade architecture (with tools like Kafka, FastAPI, etc.)**.
