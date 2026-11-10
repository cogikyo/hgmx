import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const letters = [el.querySelector('path[d^="M10 9.3"]'), el.querySelector('path[d^="M17 9.3"]')];
  animate(letters, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.12, { startDelay: τmin }) });
}
