--
-- PostgreSQL database dump
--

-- Dumped from database version 14.15 (Ubuntu 14.15-0ubuntu0.22.04.1)
-- Dumped by pg_dump version 14.15 (Ubuntu 14.15-0ubuntu0.22.04.1)

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

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


--
-- Name: process_cart_on_order_success(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.process_cart_on_order_success() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF NEW.status = 'completed' THEN
        -- Giảm số lượng stock của từng sản phẩm trong order
        UPDATE products
        SET quantity = products.quantity - od.quantity
        FROM order_detail od
        WHERE od.product_id = products.id
          AND od.order_id = NEW.id;

        -- Cập nhật số lượng trong giỏ hàng user khác
        UPDATE cart
        SET quantity = LEAST(cart.quantity, products.quantity)
        FROM products
        WHERE cart.product_id = products.id;

        -- Xóa sản phẩm hết hàng khỏi giỏ hàng
        DELETE FROM cart
        WHERE product_id IN (
            SELECT id FROM products WHERE quantity = 0
        );
    END IF;

    RETURN NEW;
END;
$$;


ALTER FUNCTION public.process_cart_on_order_success() OWNER TO postgres;

--
-- Name: update_cart_timestamp(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_cart_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
	NEW.updated_at = CURRENT_TIMESTAMP;
	RETURN NEW;
END
$$;


ALTER FUNCTION public.update_cart_timestamp() OWNER TO postgres;

--
-- Name: update_order_timestamp(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_order_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
	NEW.updated_at = CURRENT_TIMESTAMP;
	RETURN NEW;
END
$$;


ALTER FUNCTION public.update_order_timestamp() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: cart; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cart (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    product_id uuid NOT NULL,
    quantity integer NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.cart OWNER TO postgres;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    name character varying(255) NOT NULL
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.categories_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.categories_id_seq OWNER TO postgres;

--
-- Name: categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.categories_id_seq OWNED BY public.categories.id;


--
-- Name: order_detail; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_detail (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    order_id uuid,
    product_id uuid,
    quantity integer NOT NULL,
    price integer NOT NULL
);


ALTER TABLE public.order_detail OWNER TO postgres;

--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    name character varying(255),
    total_price integer NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying,
    phone character varying(20) NOT NULL,
    address text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    price integer NOT NULL,
    quantity integer DEFAULT 0,
    category_id integer NOT NULL,
    description text,
    image_url character varying(255)
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(255) NOT NULL,
    name character varying(255),
    phone character varying(20),
    address text,
    role character varying(50) NOT NULL,
    avatar character varying(255),
    CONSTRAINT users_role_check CHECK (((role)::text = ANY ((ARRAY['Admin'::character varying, 'Customer'::character varying])::text[])))
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.categories_id_seq'::regclass);


--
-- Data for Name: cart; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cart (id, user_id, product_id, quantity, updated_at) FROM stdin;
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, name) FROM stdin;
1	Laptop
2	Keyboard
3	Mouse
\.


--
-- Data for Name: order_detail; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_detail (id, order_id, product_id, quantity, price) FROM stdin;
2b913ff3-3ded-41f2-afb6-d638fc8dd27d	afc4b881-9dda-4671-8e09-ccb2de89270b	43cad845-b5eb-41d4-b28a-de4d3c1d46c4	2	20490000
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (id, user_id, name, total_price, status, phone, address, created_at, updated_at) FROM stdin;
afc4b881-9dda-4671-8e09-ccb2de89270b	66a5ca43-20e5-46b1-bf58-af92c98702e5	Đoàn Tiến Dũng	40980000	pending	0983039718	242 Xã Đàn, Phường Phương Liên, Quận Đống Đa	2024-12-31 18:16:09.823377	2024-12-31 18:16:09.823377
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, name, price, quantity, category_id, description, image_url) FROM stdin;
43cad845-b5eb-41d4-b28a-de4d3c1d46c4	Laptop ASUS TUF Gaming F15 FX507ZC4-HN095W	20490000	10	1	\N	Lap1-ASUS-TUF-F15.png
ed7e4c9a-9025-4dae-b1a8-942eda8c0e16	Laptop Lenovo LOQ 15ARP9 83JC007HVN	20990000	10	1	\N	Lap2-Lenovo-LOQ.png
3a08aaa2-412d-4404-a41e-22e73c2e2e84	Laptop Dell Inspiron 15 3520	16390000	10	1	\N	Lap3-DELL-inspiron-15.png
c467f401-aa9d-4473-a78a-91cc69ef6b95	Logitech G304 Lightspeed	745000	20	3	\N	Mouse4-Logitech-G304.png
daf870b0-24f5-4f4a-91e6-7bcc7c40c793	E-DRA-EK375-Alpha	1199000	30	2	\N	Keyboard5-E-DRA-EK375.png
f094a2b0-2251-4a5e-82d2-5096b9ad7caa	AULA-F75	490000	30	2	\N	Keyboard6-AULA-F75.png
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, password, name, phone, address, role, avatar) FROM stdin;
73ce2f7c-ca34-4feb-8fc7-d8e9978cf192	sangpham1224@gmail.com	$2a$10$WOyBEFq4mv7DEzH6pAWYk.bF8E9ocSy39zsYB67/4eVPf2P91a9Ya	Phạm Trường Sang	\N	\N	Customer	\N
66a5ca43-20e5-46b1-bf58-af92c98702e5	shiraishi2612@gmail.com	$2a$10$OTgyyXhaRfk.tnwaN20tO.V6rP2Hdumj2HMFYU4iLWJKM85QJ5ebe	Đoàn Tiến Dũng	0983039718	242 Xã Đàn, Phường Phương Liên, Quận Đống Đa	Customer	/static/media/img.9403b305df778ea51be8.png
\.


--
-- Name: categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.categories_id_seq', 3, true);


--
-- Name: cart cart_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart
    ADD CONSTRAINT cart_pkey PRIMARY KEY (id);


--
-- Name: cart cart_user_id_product_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart
    ADD CONSTRAINT cart_user_id_product_id_key UNIQUE (user_id, product_id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: order_detail order_detail_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_detail
    ADD CONSTRAINT order_detail_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: orders order_status_update; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER order_status_update AFTER UPDATE ON public.orders FOR EACH ROW WHEN ((((new.status)::text = 'completed'::text) AND ((old.status)::text <> 'completed'::text))) EXECUTE FUNCTION public.process_cart_on_order_success();


--
-- Name: cart update_cart_before_update; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_cart_before_update BEFORE UPDATE ON public.cart FOR EACH ROW EXECUTE FUNCTION public.update_cart_timestamp();


--
-- Name: orders update_order_before_update; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_order_before_update BEFORE UPDATE ON public.orders FOR EACH ROW EXECUTE FUNCTION public.update_order_timestamp();


--
-- Name: cart cart_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart
    ADD CONSTRAINT cart_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: cart cart_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart
    ADD CONSTRAINT cart_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: order_detail order_detail_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_detail
    ADD CONSTRAINT order_detail_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: order_detail order_detail_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_detail
    ADD CONSTRAINT order_detail_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

