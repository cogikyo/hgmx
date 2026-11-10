import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector('path[d^="M21 11.5"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M16 5h6"], path[d="M19 2v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });

  const sun = el.querySelector("circle");
  animate(sun, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
