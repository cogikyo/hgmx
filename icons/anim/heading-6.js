import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const stems = [el.querySelector('path[d="M4 18V6"]'), el.querySelector('path[d="M12 18V6"]')];
  animate(stems, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin / 2) });
  const bar = el.querySelector('path[d="M4 12h8"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const numeral = el.querySelectorAll('path:not([d$="V6"], [d="M4 12h8"]), circle');
  animate(numeral, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
