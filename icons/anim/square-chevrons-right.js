import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const chevrons = [el.querySelector('path[d^="M14 15"]'), el.querySelector('path[d^="M8 15"]')];
  animate(chevrons, { translateX: [0, 2, -1, 0] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
