import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M3 3"]');
  animate(axes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const links = el.querySelectorAll('path:not([d^="M3 3"])');
  animate(links, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const nodes = [el.querySelector('circle[cx="12"]'), el.querySelector('circle[cx="16"]'), el.querySelector('circle[cx="9"]')];
  animate(nodes, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
