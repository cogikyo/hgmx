import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector('path[d^="M20 16"]');
  animate(head, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const antenna = el.querySelector('path[d^="M12 6V2"]');
  animate(antenna, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const eyes = [el.querySelector('path[d="M9 11v2"]'), el.querySelector('path[d="M15 11v2"]')];
  animate(eyes, { translateY: [0, -2, 1, 0] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
