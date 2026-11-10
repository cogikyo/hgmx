import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('path[d^="M8 6"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const drivers = [el.querySelector('path[d="M17 9h.01"]'), el.querySelector("circle")];
  animate(drivers, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });
}
