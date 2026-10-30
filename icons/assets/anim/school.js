import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const walls = el.querySelectorAll('path[d^="m6 11"], path[d^="M18 4.933"], path[d^="M6 4.933"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const roof = el.querySelector('path[d^="m4 6"]');
  animate(roof, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const clock = el.querySelector("circle");
  animate(clock, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
