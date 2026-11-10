import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const penguin = el.querySelector("g");
  animate(penguin, { rotate: [0, -8, 8, -4, 0] }, { duration: τmax });

  const left = el.querySelector('path[d^="M9 10c"]');
  animate(left, { rotate: [0, 20, -5, 0] }, { duration: τmid, delay: τmin });
  const right = el.querySelector('path[d^="M15 10c"]');
  animate(right, { rotate: [0, -20, 5, 0] }, { duration: τmid, delay: τmin });
}
