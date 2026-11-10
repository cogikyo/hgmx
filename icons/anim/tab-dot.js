import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tab = el.querySelector('path[d^="M4 20"]');
  animate(tab, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [1, 1.75, 1] }, { duration: τmid, delay: τmin });
}
