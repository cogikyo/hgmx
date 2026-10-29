import { animate, τmax } from "../motion.js";

export default function (el) {
  const fork = el.querySelector("g");
  animate(fork, { rotate: [0, 12, -4, 0] }, { duration: τmax });

  const knife = el.querySelector('path[d^="M21 15"]');
  animate(knife, { rotate: [0, -12, 4, 0] }, { duration: τmax });
}
