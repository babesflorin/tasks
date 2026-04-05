<?php

namespace App\Tests\DataFixtures;

use App\Domain\Entity\Task;
use Doctrine\Bundle\FixturesBundle\Fixture;
use Doctrine\Persistence\ObjectManager;

class TasksFixtures extends Fixture
{
    public function load(ObjectManager $manager)
    {
        for ($i = 0; $i < 5; $i++) {
            $when = (new \DateTime())->modify("+$i Day");
            $task = new Task("Task name $i", "Task description $i", $when);
            $manager->persist($task);
        }
        $manager->flush();
    }
}
