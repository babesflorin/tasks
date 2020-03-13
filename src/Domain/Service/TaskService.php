<?php

namespace App\Domain\Service;

use App\Domain\Dto\TaskCollectionDto;
use App\Domain\Dto\TaskDto;
use App\Domain\Exception\InvalidTaskException;
use App\Domain\Exception\TaskNotFoundException;
use App\Domain\Repository\TaskRepositoryInterface;
use App\Domain\Transformer\TaskCollectionDtoToCollectionTransformer;
use App\Domain\Transformer\TaskDtoToEntityTransformer;
use App\Domain\Validator\TaskValidator;

class TaskService
{
    /**
     * @var TaskRepositoryInterface
     */
    private $repository;
    /**
     * @var TaskValidator
     */
    private $validator;
    /**
     * @var TaskDtoToEntityTransformer
     */
    private $dtoTransformer;
    /**
     * @var TaskCollectionDtoToCollectionTransformer
     */
    private $dtoCollectionTransformer;

    public function __construct(
        TaskRepositoryInterface $repository,
        TaskValidator $validator,
        TaskDtoToEntityTransformer $dtoTransformer,
        TaskCollectionDtoToCollectionTransformer $dtoCollectionTransformer
    ) {
        $this->repository = $repository;
        $this->validator = $validator;
        $this->dtoTransformer = $dtoTransformer;
        $this->dtoCollectionTransformer = $dtoCollectionTransformer;
    }

    /**
     * @throws InvalidTaskException
     */
    public function addTask(TaskDto $taskDto): TaskDto
    {
        $this->validator->validate($taskDto);
        $task = $this->repository->saveTask($this->dtoTransformer->transform($taskDto));

        return $this->dtoTransformer->reverseTransform($task);
    }

    public function getAllTasks(): TaskCollectionDto
    {
        return $this->dtoCollectionTransformer->reverseTransform($this->repository->getTasks());
    }

    public function completeTask(int $taskId): TaskDto
    {
        $task = $this->repository->findTaskById($taskId);
        if (null === $task) {
            throw new TaskNotFoundException("Task not found!");
        }
        $task->complete();
        $this->repository->saveTask($task);

        return $this->dtoTransformer->reverseTransform($task);
    }

    /**
     * @throws InvalidTaskException
     * @throws TaskNotFoundException
     */
    public function updateTask(TaskDto $taskDto) : TaskDto
    {
        $this->validator->validate($taskDto, true);
        $task = $this->repository->findTaskById($taskDto->id);
        if (null === $task) {
            throw new TaskNotFoundException("Task not found!");
        }
        $task = $this->dtoTransformer->transform($taskDto, $task);
        $this->repository->saveTask($task);

        return $this->dtoTransformer->reverseTransform($task);
    }

    public function getTask(int $taskId): TaskDto
    {
        $task = $this->repository->findTaskById($taskId);
        if (null === $task) {
            throw new TaskNotFoundException("Task not found!");
        }

        return $this->dtoTransformer->reverseTransform($task);
    }
}
