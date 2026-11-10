import { animate, τmax } from "../motion.js";

export default function (el) {
  const fork = el.querySelector("g");
  animate(fork, { rotate: [0, 12, -4, 0] }, { duration: τmax });

  const spoon = el.querySelector('path[d^="M15 15 3.3"]');
  animate(spoon, { rotate: [0, -12, 4, 0] }, { duration: τmax });
}
