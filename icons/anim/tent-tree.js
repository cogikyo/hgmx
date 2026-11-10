import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tent = el.querySelectorAll('path[d^="M17 14H7"], path[d="M8 14v8"], path[d="m9 14 5 8"]');
  animate(tent, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tree = el.querySelector("g");
  animate(tree, { rotate: [0, -10, 6, 0] }, { duration: τmid, delay: τmin });
}
