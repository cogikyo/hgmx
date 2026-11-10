import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cable = el.querySelector('path[d^="m22 11"]');
  animate(cable, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const hangers = el.querySelectorAll('path[d$="9.728V16"]');
  animate(hangers, { scaleY: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
}
