--
-- PostgreSQL database dump
--

-- Dumped from database version 14.5 (Debian 14.5-1.pgdg110+1)
-- Dumped by pg_dump version 14.5 (Homebrew)

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



DROP TABLE [IF EXIST] genres, movies, movies_genres, users [CASCADE | RESTRICT];

--
-- Name: products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.products (
    id integer NOT NULL,
    product_name character varying(255),
    price integer NOT NULL,
    description character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.products ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.products_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: stocks; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.stocks (
    id integer NOT NULL,
    product_id character varying(512),
    quantity integer
);


--
-- Name: stocks_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.stocks ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.stocks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: histories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.histories (
    id integer NOT NULL,
    user_id integer,
    product_id integer,
    quantity integer,
    discount_applied numeric
);


--
-- Name: coupons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE IF not exists coupons (
                                       id INT NOT NULL,
                                       code VARCHAR(50),
                                       rate INT NOT NULL,
                                       product_id INT NOT NULL,
                                       user_id INT,
                                       active BOOLEAN DEFAULT true,
                                       expired_at TIMESTAMP,
                                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       PRIMARY KEY (id),
                                       FOREIGN KEY (product_id) REFERENCES products(id),
                                       FOREIGN KEY (user_id) REFERENCES users(id)
)


--
-- Name: histories_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.histories ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.histories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id integer NOT NULL,
    first_name character varying(255),
    last_name character varying(255),
    email character varying(255),
    password character varying(255),
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.users ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.products (id, product_name, price, description, created_at, updated_at) FROM stdin;
1	pen 5 study	2022-09-23 00:00:00	2022-09-23 00:00:00
2	book 10 study	2022-09-23 00:00:00	2022-09-23 00:00:00
3	ruler 2 study	2022-09-23 00:00:00	2022-09-23 00:00:00
4	backpack 50 study	2022-09-23 00:00:00	2022-09-23 00:00:00
5	cap 15 wear	2022-09-23 00:00:00	2022-09-23 00:00:00
\.


--
-- Data for Name: histories; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.histories (id, user_id, product_id, quantity, discount, created_at, updated_at) FROM stdin;
1	1   1   5   0.2      2022-09-23 00:00:00	2022-09-23 00:00:00
2	1   3   3   0.2      2022-09-23 00:00:00	2022-09-23 00:00:00
3	1   4   10   0.2      2022-09-23 00:00:00	2022-09-23 00:00:00
\.


--
-- Data for Name: stocks; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.stocks (id, product_id, quantity) FROM stdin;
1	1	5
2	2	12
3	3	5
4	4	11
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, first_name, last_name, email, password, created_at, updated_at) FROM stdin;
1	Admin	User	admin@example.com	$2a$14$wVsaPvJnJJsomWArouWCtusem6S/.Gauq/GjOIEHpyh2DAMmso1wy	2022-09-23 00:00:00	2022-09-23 00:00:00
2	Danh	Huynh	danh@example.com	$2a$14$wVsaPvJnJJsomWArouWCtusem6S/.Gauq/GjOIEHpyh2DAMmso1wy	2022-09-23 00:00:00	2022-09-23 00:00:00
\.


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.products_id_seq', 5, true);


--
-- Name: stocks_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.stocks_id_seq', 4, true);


--
-- Name: histories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.histories_id_seq', 3, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.users_id_seq', 2, true);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: histories histories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.histories
    ADD CONSTRAINT histories_pkey PRIMARY KEY (id);


--
-- Name: stocks stocks_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.stocks
    ADD CONSTRAINT stocks_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: histories histories_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.histories
    ADD CONSTRAINT histories_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: histories histories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.histories
    ADD CONSTRAINT histories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

