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
    graph_id CHARACTER(36) NOT NULL,
    node_id INTEGER NOT NULL,
    shape CHARACTER VARYING(32),
    maxheight INTEGER,
    minheight INTEGER
);


-- edges
CREATE TABLE edges (
    graph_id CHARACTER(36) NOT NULL,
    parent_id INTEGER NOT NULL,
    child_id INTEGER NOT NULL,
    label CHARACTER VARYING(64) NOT NULL,
    shaper CHARACTER VARYING(512),
    coorder CHARACTER VARYING(512)
);


-- labels assigned to nodes
CREATE TABLE node_tags (
    graph_id CHARACTER(36) NOT NULL,
    node_id INTEGER NOT NULL,
    tag_key CHARACTER VARYING(64) NOT NULL,
    tag_val CHARACTER VARYING(64)
);


-- data assigned to nodes
CREATE TABLE node_data (
    graph_id CHARACTER(36) NOT NULL,
    node_id INTEGER NOT NULL,
    data DECIMAL[]
);


-- primary keys
ALTER TABLE ONLY nodes
    ADD CONSTRAINT pk_nodes PRIMARY KEY (graph_id, node_id);

ALTER TABLE ONLY edges
    ADD CONSTRAINT pk_edges PRIMARY KEY (graph_id, parent_id, child_id, label);

ALTER TABLE ONLY node_tags
    ADD CONSTRAINT pk_node_labels PRIMARY KEY (graph_id, node_id, tag_key);

ALTER TABLE ONLY node_data
    ADD CONSTRAINT pk_node_data PRIMARY KEY (graph_id, node_id);


-- indexing
CREATE UNIQUE INDEX index_nodes_graph_node_id
    ON nodes USING BTREE (graph_id, node_id);

CREATE INDEX index_edges_parent_id
    ON edges USING BTREE (graph_id, parent_id);

CREATE INDEX index_edges_child_id
    ON edges USING BTREE (graph_id, child_id);

CREATE INDEX index_node_labels_id ON node_tags USING BTREE (graph_id, node_id);

CREATE INDEX index_node_data_id ON node_data USING BTREE (graph_id, node_id);


-- foreign keys
ALTER TABLE ONLY edges ADD CONSTRAINT fk_edges_graph_parent_id
    FOREIGN KEY (graph_id, parent_id) REFERENCES nodes(graph_id, node_id);

ALTER TABLE ONLY edges ADD CONSTRAINT fk_edges_graph_child_id
    FOREIGN KEY (graph_id, child_id) REFERENCES nodes(graph_id, node_id);

ALTER TABLE ONLY node_tags ADD CONSTRAINT fk_node_graph_node_id
    FOREIGN KEY (graph_id, node_id) REFERENCES nodes(graph_id, node_id);

ALTER TABLE ONLY node_data ADD CONSTRAINT fk_node_data_id
    FOREIGN KEY (graph_id, node_id) REFERENCES nodes(graph_id, node_id);
