import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const brain = el.querySelector("g.brain");
  animate(brain, { scale: [1, 0.9, 1.1, 1] }, { duration: τmid });

  const gear = el.querySelector("g.gear");
  animate(gear, { rotate: [0, 360], scale: [1, 0.8, 1] }, { duration: τmid, delay: τmin });
}
