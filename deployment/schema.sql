--
-- PostgreSQL database dump
--

-- Dumped from database version 15.13 (Debian 15.13-1.pgdg120+1)
-- Dumped by pg_dump version 15.13 (Debian 15.13-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: analysis_packages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.analysis_packages (
    id uuid NOT NULL,
    name text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    "interval" text NOT NULL,
    time_from timestamp with time zone NOT NULL,
    time_to timestamp with time zone NOT NULL,
    hist_bins integer NOT NULL,
    hist_min double precision NOT NULL,
    hist_max double precision NOT NULL,
    mcap_min bigint,
    inception_max timestamp with time zone,
    symbol_count integer,
    status text NOT NULL,
    user_id uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid NOT NULL
);


--
-- Name: analysis_results; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.analysis_results (
    package_id uuid NOT NULL,
    ticker text NOT NULL,
    count integer NOT NULL,
    mean double precision NOT NULL,
    stddev double precision NOT NULL,
    variance double precision NOT NULL,
    min double precision NOT NULL,
    max double precision NOT NULL,
    histogram jsonb NOT NULL,
    chart_path text
);


--
-- Name: batch_update_log; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.batch_update_log (
    id integer NOT NULL,
    updater_name character varying(100) NOT NULL,
    started_at timestamp without time zone NOT NULL,
    completed_at timestamp without time zone,
    status character varying(20) NOT NULL,
    symbols_processed integer DEFAULT 0,
    symbols_updated integer DEFAULT 0,
    error_message text
);


--
-- Name: batch_update_log_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.batch_update_log_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: batch_update_log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.batch_update_log_id_seq OWNED BY public.batch_update_log.id;


--
-- Name: errors; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.errors (
    id integer NOT NULL,
    "timestamp" timestamp with time zone DEFAULT now(),
    source text NOT NULL,
    error_type text NOT NULL,
    message text NOT NULL,
    details text
);


--
-- Name: errors_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.errors_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: errors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.errors_id_seq OWNED BY public.errors.id;


--
-- Name: monthly_prices; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.monthly_prices (
    date timestamp with time zone NOT NULL,
    close double precision,
    symbol_ticker text NOT NULL,
    open double precision DEFAULT '0'::double precision,
    high double precision DEFAULT '0'::double precision,
    low double precision DEFAULT '0'::double precision,
    avg double precision DEFAULT '0'::double precision,
    yoy double precision
);


--
-- Name: notes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notes (
    id text NOT NULL,
    type text,
    title text,
    content text,
    date timestamp with time zone,
    data text
);


--
-- Name: symbols; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.symbols (
    ticker text NOT NULL,
    exchange text,
    name text,
    type text,
    currency text,
    sector text,
    industry text,
    country text,
    description text,
    website text,
    isin text,
    inception timestamp with time zone,
    analytics json,
    details json,
    last_price_update timestamp with time zone,
    last_profile_update timestamp with time zone,
    last_profile_status text,
    last_price_status text,
    is_actively_trading boolean,
    market_cap bigint,
    oldest_price timestamp with time zone,
    primary_listing text,
    cik text,
    ath12m double precision,
    current_price_usd double precision,
    current_price_time timestamp with time zone
);


--
-- Name: user_favorites; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_favorites (
    ticker character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    user_id uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid NOT NULL
);


--
-- Name: user_ratings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_ratings (
    id integer NOT NULL,
    ticker character varying(20) NOT NULL,
    rating integer NOT NULL,
    notes text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    user_id uuid DEFAULT '00000000-0000-0000-0000-000000000000'::uuid NOT NULL,
    CONSTRAINT user_ratings_rating_check CHECK (((rating >= '-5'::integer) AND (rating <= 5)))
);


--
-- Name: user_ratings_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.user_ratings_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: user_ratings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.user_ratings_id_seq OWNED BY public.user_ratings.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    name character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


--
-- Name: weekly_prices; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.weekly_prices (
    date timestamp with time zone NOT NULL,
    close double precision,
    symbol_ticker text NOT NULL,
    open double precision DEFAULT '0'::double precision,
    high double precision DEFAULT '0'::double precision,
    low double precision DEFAULT '0'::double precision,
    avg double precision DEFAULT '0'::double precision,
    yoy double precision
);


--
-- Name: batch_update_log id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.batch_update_log ALTER COLUMN id SET DEFAULT nextval('public.batch_update_log_id_seq'::regclass);


--
-- Name: errors id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.errors ALTER COLUMN id SET DEFAULT nextval('public.errors_id_seq'::regclass);


--
-- Name: user_ratings id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_ratings ALTER COLUMN id SET DEFAULT nextval('public.user_ratings_id_seq'::regclass);


--
-- Name: analysis_packages analysis_packages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.analysis_packages
    ADD CONSTRAINT analysis_packages_pkey PRIMARY KEY (id);


--
-- Name: analysis_results analysis_results_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.analysis_results
    ADD CONSTRAINT analysis_results_pkey PRIMARY KEY (package_id, ticker);


--
-- Name: batch_update_log batch_update_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.batch_update_log
    ADD CONSTRAINT batch_update_log_pkey PRIMARY KEY (id);


--
-- Name: errors errors_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.errors
    ADD CONSTRAINT errors_pkey PRIMARY KEY (id);


--
-- Name: symbols idx_16389_sqlite_autoindex_symbols_1; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.symbols
    ADD CONSTRAINT idx_16389_sqlite_autoindex_symbols_1 PRIMARY KEY (ticker);


--
-- Name: weekly_prices idx_16394_sqlite_autoindex_weekly_prices_1; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.weekly_prices
    ADD CONSTRAINT idx_16394_sqlite_autoindex_weekly_prices_1 PRIMARY KEY (date, symbol_ticker);


--
-- Name: monthly_prices idx_16403_sqlite_autoindex_monthly_prices_1; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.monthly_prices
    ADD CONSTRAINT idx_16403_sqlite_autoindex_monthly_prices_1 PRIMARY KEY (date, symbol_ticker);


--
-- Name: notes idx_16521_sqlite_autoindex_notes_1; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notes
    ADD CONSTRAINT idx_16521_sqlite_autoindex_notes_1 PRIMARY KEY (id);


--
-- Name: user_favorites user_favorites_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_favorites
    ADD CONSTRAINT user_favorites_pkey PRIMARY KEY (ticker);


--
-- Name: user_favorites user_favorites_user_ticker_unique; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_favorites
    ADD CONSTRAINT user_favorites_user_ticker_unique UNIQUE (user_id, ticker);


--
-- Name: user_ratings user_ratings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.user_ratings
    ADD CONSTRAINT user_ratings_pkey PRIMARY KEY (user_id, ticker, created_at);


--
-- Name: users users_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_name_key UNIQUE (name);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_16394_idx_weekly_symbol_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_16394_idx_weekly_symbol_date ON public.weekly_prices USING btree (symbol_ticker, date);


--
-- Name: idx_16403_idx_monthly_symbol_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_16403_idx_monthly_symbol_date ON public.monthly_prices USING btree (symbol_ticker, date);


--
-- Name: idx_analysis_mean; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_analysis_mean ON public.analysis_results USING btree (package_id, mean);


--
-- Name: idx_analysis_packages_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_analysis_packages_user ON public.analysis_packages USING btree (user_id);


--
-- Name: idx_analysis_stddev; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_analysis_stddev ON public.analysis_results USING btree (package_id, stddev);


--
-- Name: idx_analysis_variance; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_analysis_variance ON public.analysis_results USING btree (package_id, variance);


--
-- Name: idx_errors_source; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_errors_source ON public.errors USING btree (source);


--
-- Name: idx_errors_timestamp; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_errors_timestamp ON public.errors USING btree ("timestamp" DESC);


--
-- Name: idx_user_favorites_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_favorites_user ON public.user_favorites USING btree (user_id);


--
-- Name: idx_user_ratings_ticker; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_ratings_ticker ON public.user_ratings USING btree (ticker);


--
-- Name: idx_user_ratings_ticker_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_ratings_ticker_created ON public.user_ratings USING btree (ticker, created_at DESC);


--
-- Name: idx_user_ratings_user; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_ratings_user ON public.user_ratings USING btree (user_id);


--
-- Name: idx_user_ratings_user_ticker; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_user_ratings_user_ticker ON public.user_ratings USING btree (user_id, ticker);


--
-- Name: analysis_results analysis_results_package_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.analysis_results
    ADD CONSTRAINT analysis_results_package_id_fkey FOREIGN KEY (package_id) REFERENCES public.analysis_packages(id) ON DELETE CASCADE;


--
-- Name: monthly_prices monthly_prices_symbol_ticker_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.monthly_prices
    ADD CONSTRAINT monthly_prices_symbol_ticker_fkey FOREIGN KEY (symbol_ticker) REFERENCES public.symbols(ticker);


--
-- Name: weekly_prices weekly_prices_symbol_ticker_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.weekly_prices
    ADD CONSTRAINT weekly_prices_symbol_ticker_fkey FOREIGN KEY (symbol_ticker) REFERENCES public.symbols(ticker);


--
-- PostgreSQL database dump complete
--

