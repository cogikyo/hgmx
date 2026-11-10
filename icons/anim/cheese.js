import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const wedge = el.querySelectorAll('path[d^="M21 19"], path[d="M9 12H3"], path[d^="M9 12c"], path[d="M21 12h-6"]');
  animate(wedge, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hole = el.querySelector("circle");
  animate(hole, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
