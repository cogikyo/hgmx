import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M21 15.251"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const arc = el.querySelector('path[d^="M8 19"]');
  animate(arc, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const head = el.querySelector('path[d^="M7 11"]');
  animate(head, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
