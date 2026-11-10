import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const lines = el.querySelectorAll('path[d$="h11"]');
  animate(lines, { translateX: [-3, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06) });

  const chevrons = [el.querySelector('path[d^="m3 10 "]'), el.querySelector('path[d^="m3 20 "]')];
  animate(chevrons, { translateX: [0, 2, -0.5, 0] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
