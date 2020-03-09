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

    public function getTasks(): TaskCollection
    {
        $tasks = $this->repository->findAll();

        return new TaskCollection($tasks);
    }
}