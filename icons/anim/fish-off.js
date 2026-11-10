import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const fish = el.querySelector("g");
  animate(fish, { rotate: [0, -6, 4, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m16.01"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
