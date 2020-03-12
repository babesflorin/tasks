<?php

namespace App\Domain\Transformer;

use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;

class TaskDtoToEntityTransformer
{
    public function transform(TaskDto $taskDto, ?Task $task = null)
    {
        $when = \DateTime::createFromFormat('Y-m-d', $taskDto->when)->setTime(0, 0, 0);
        if (null === $task) {
            $task = new Task(
                $taskDto->name,
                $taskDto->description,
                $when
            );
        } else {
            $task->updateFromRaw(
                $taskDto->name,
                $taskDto->description,
                $when
            );
        }

        return $task;
    }

    public function reverseTransform(Task $task)
    {
        $taskDto = new TaskDto();
        $taskDto->id = $task->getId();
        $taskDto->name = $task->getName();
        $taskDto->description = $task->getDescription();
        $taskDto->when = $task->getWhen()->format('Y-m-d');
        $taskDto->done = $task->isDone();
        $taskDto->createdAt = $task->getCreatedAt();
        $taskDto->updatedAt = $task->getUpdateAt();

        return $taskDto;
    }
}
