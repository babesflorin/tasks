<?php

namespace App\Tests\Domain\Entity;

use App\Domain\Entity\Task;
use App\Tests\TestCase;

class TaskTest extends TestCase
{

    public function testComplete()
    {
        $task = new Task("name", "description", new \DateTime());
        $this->assertFalse($task->isDone());
        $task->complete();
        $this->assertTrue($task->isDone());
    }
}
