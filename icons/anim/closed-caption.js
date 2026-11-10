import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const letters = [el.querySelector('path[d^="M10 9.17"]'), el.querySelector('path[d^="M17 9.17"]')];
  animate(letters, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.12, { startDelay: τmin }) });
}
