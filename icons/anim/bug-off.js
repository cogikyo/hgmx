import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bug = el.querySelector("g");
  animate(bug, { rotate: [0, -7, 4, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
