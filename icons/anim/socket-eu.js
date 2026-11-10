import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const plate = el.querySelector("rect");
  animate(plate, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pins = [el.querySelector('path[d="M10 12h.01"]'), el.querySelector('path[d="M14 12h.01"]')];
  animate(pins, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
