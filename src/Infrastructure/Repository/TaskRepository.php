<?php

namespace App\Infrastructure\Repository;

use App\Domain\Entity\Task;
use App\Domain\Entity\TaskCollection;
use App\Domain\Repository\TaskRepositoryInterface;
use Doctrine\ORM\EntityManagerInterface;
use Doctrine\ORM\EntityRepository;

class TaskRepository implements TaskRepositoryInterface
{
    /** @var EntityManagerInterface */
    private $entityManager;
    /** @var EntityRepository */
    private $repository;

    public function __construct(EntityManagerInterface $entityManager)
    {
        $this->entityManager = $entityManager;
        $this->repository = $entityManager->getRepository(Task::class);
    }

    public function saveTask(Task $task): Task
    {
        $this->entityManager->persist($task);
        $this->entityManager->flush();

        return $task;
    }

    public function getTasks(?bool $areDone = null, ?\DateTime $when = null): TaskCollection
    {
        $criteria =[];

        if (null !== $areDone) {
            $criteria['done'] = $areDone;
        }
        if (null !== $when) {
            $criteria['when'] = $when->setTime(0, 0, 0);
        }
        $tasks = $this->repository->findBy($criteria);

        return new TaskCollection($tasks);
    }

    public function findTaskById(int $taskId): ?Task
    {
        return $this->repository->find($taskId);
    }

    public function deleteTask(Task $task): bool
    {
        $this->entityManager->remove($task);
        $this->entityManager->flush();

        return null === $task->getId();
    }
}
