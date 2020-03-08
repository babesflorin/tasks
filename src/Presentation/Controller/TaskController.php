<?php

namespace App\Presentation\Controller;

use App\Domain\Dto\TaskDto;
use App\Domain\Service\TaskService;
use App\Presentation\Transformer\TaskTransformer;
use League\Fractal\Manager;
use League\Fractal\Resource\Item;
use Sensio\Bundle\FrameworkExtraBundle\Configuration\ParamConverter;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\Routing\Annotation\Route;

/**
 * @Route("/task")
 */
class TaskController extends AbstractController
{
    /**
     * @var TaskService
     */
    private $service;

    /**
     * @var Manager
     */
    private $fractal;

    /**
     * @var TaskTransformer
     */
    private $taskTransformer;

    public function __construct(TaskService $service, Manager $fractal, TaskTransformer $taskTransformer)
    {
        $this->service = $service;
        $this->fractal = $fractal;
        $this->taskTransformer = $taskTransformer;
    }

    /**
     * @Route(
     *     "",
     *     name="task_add",
     *     methods={"POST"},
     *     format="application/json",
     *     requirements={
     *          "_format" : "application/json"
     *      }
     * )
     * @ParamConverter("task", class=TaskDto::class)
     */
    public function index(TaskDto $task)
    {
        $taskDto = $this->service->addTask($task);
        $resource = new Item($taskDto, $this->taskTransformer);

        return $this->json(
            $this->fractal->createData($resource)->toArray()
        );
    }
}
