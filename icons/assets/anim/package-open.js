import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelectorAll('path[d^="M20 13"], path[d^="M12 22"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flaps = el.querySelectorAll('path[d^="M15.17"], path[d^="M21 12.43"]');
  animate(flaps, { translateY: [0, -3, 1, 0] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
