import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const chevrons = [el.querySelector('path[d^="M14 15"]'), el.querySelector('path[d^="M8 15"]')];
  animate(chevrons, { translateX: [0, 1.5, -0.5, 0] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
