import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const virus = el.querySelector("g");
  animate(virus, { rotate: [0, 20, 0], scale: [1, 0.9, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
