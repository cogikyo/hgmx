import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const base = el.querySelector("line");
  animate(base, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const dashes = [el.querySelector('path[d="M18 6V5"]'), el.querySelector('path[d="M18 11v-1"]')];
  animate(dashes, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09, { startDelay: τmin }) });

  const head = el.querySelector('circle[cx="18"]');
  animate(head, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
