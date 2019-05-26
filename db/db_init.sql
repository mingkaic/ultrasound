SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET default_tablespace = '';
SET default_with_oids = false;

-- nodes
CREATE TABLE nodes (
    id INTEGER NOT NULL,
    graph_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    representation CHARACTER VARYING(64),
    shape INTEGER[8]
);

CREATE SEQUENCE nodes_id_seq
    START WITH 1 INCREMENT BY 1
    NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE nodes_id_seq OWNED BY nodes.id;


-- edges
CREATE TABLE edges (
    graph_id INTEGER NOT NULL,
    parent_id INTEGER NOT NULL,
    child_id INTEGER NOT NULL,
    label CHARACTER VARYING(64)
);


-- labels assigned to nodes
CREATE TABLE node_labels (
    graph_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    label CHARACTER VARYING(64)
);


-- data assigned to nodes
CREATE TABLE node_data (
    id INTEGER NOT NULL,
    index INTEGER NOT NULL,
    datum DECIMAL NOT NULL
);


-- primary keys
ALTER TABLE ONLY nodes
    ADD CONSTRAINT pk_nodes PRIMARY KEY (id);

ALTER TABLE ONLY edges
    ADD CONSTRAINT pk_edges PRIMARY KEY (graph_id, parent_id, child_id);

ALTER TABLE ONLY node_labels
    ADD CONSTRAINT pk_node_labels PRIMARY KEY (graph_id, node_id);

ALTER TABLE ONLY node_data
    ADD CONSTRAINT pk_node_data PRIMARY KEY (id, index);


-- foreign keys
CREATE UNIQUE INDEX index_nodes_graph_node_id
    ON nodes USING BTREE (graph_id, node_id);

CREATE INDEX index_edges_parent_id
    ON edges USING BTREE (graph_id, parent_id);

CREATE INDEX index_edges_child_id
    ON edges USING BTREE (graph_id, child_id);

CREATE INDEX index_node_labels_id ON node_labels USING BTREE (graph_id, node_id);

CREATE INDEX index_node_data_id ON node_data USING BTREE (id);

ALTER TABLE ONLY edges ADD CONSTRAINT fk_edges_graph_parent_id
    FOREIGN KEY (graph_id, parent_id) REFERENCES nodes(graph_id, node_id);

ALTER TABLE ONLY edges ADD CONSTRAINT fk_edges_graph_child_id
    FOREIGN KEY (graph_id, child_id) REFERENCES nodes(graph_id, node_id);

ALTER TABLE ONLY node_labels ADD CONSTRAINT fk_node_graph_node_id
    FOREIGN KEY (graph_id, node_id) REFERENCES nodes(graph_id, node_id);

ALTER TABLE ONLY node_data ADD CONSTRAINT fk_node_data_id
    FOREIGN KEY (id) REFERENCES nodes(id);
