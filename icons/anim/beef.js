import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const beef = el.querySelectorAll('path[d^="M16.4"], path[d^="m18.5"]');
  animate(beef, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bone = el.querySelector("circle");
  animate(bone, { pathLength: [0, 1], opacity: [0, 1], rotate: [270, 90] }, { duration: τmid, delay: τmin });
}
