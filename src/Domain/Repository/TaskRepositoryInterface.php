<?php

namespace App\Domain\Repository;

use App\Domain\Entity\Task;
use App\Domain\Entity\TaskCollection;

interface TaskRepositoryInterface
{
    public function saveTask(Task $task): Task;

    public function getTasks(): TaskCollection;

    public function findTaskById(int $taskId): ?Task;
}
