CREATE TABLE `quizzes`
(
    `id`         integer,
    `created_at` datetime,
    `updated_at` datetime,
    `deleted_at` datetime,
    `question`   text,
    PRIMARY KEY (`id`)
);

CREATE TABLE `answers`
(
    `id`         integer,
    `created_at` datetime,
    `updated_at` datetime,
    `deleted_at` datetime,
    `text`       text,
    `is_correct`  bool,
    `quiz_id`    integer,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_quizzes_answers` FOREIGN KEY (`quiz_id`) REFERENCES `quizzes` (`id`)
);

CREATE TABLE `games`
(
    `id`                    integer,
    `created_at`            datetime,
    `updated_at`            datetime,
    `deleted_at`            datetime,
    `total_questions`       integer,
    `failed_questions`      integer,
    `not_answered_questions` integer,
    `start_time`            datetime,
    `end_time`              datetime,
    PRIMARY KEY (`id`)
);

INSERT INTO quizzes (id, created_at, updated_at, deleted_at, question)
VALUES (null, null, null, null, 'Что делал слон когда пришел на поле он?')
