import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const barrel = el.querySelector('path[d^="M6.23"]');
  animate(barrel, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hoops = [el.querySelector('path[d^="M4.67 7"]'), el.querySelector('path[d^="M4.67 17"]')];
  animate(hoops, { scaleX: [0, 1] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
