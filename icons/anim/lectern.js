import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const desk = el.querySelector('path[d^="M15 13"]');
  animate(desk, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const mic = el.querySelector("path[transform-origin]");
  animate(mic, { rotate: [0, -20, 10, 0] }, { duration: τmid, delay: τmin });
}
