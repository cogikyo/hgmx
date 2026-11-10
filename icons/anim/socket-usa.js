import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const plate = el.querySelector("rect");
  animate(plate, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slots = [el.querySelector('path[d="M10 11v2"]'), el.querySelector('path[d="M14 11v2"]')];
  animate(slots, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
