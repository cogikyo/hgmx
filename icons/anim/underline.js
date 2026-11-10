import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const letter = el.querySelector("path");
  animate(letter, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const rule = el.querySelector("line");
  animate(rule, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
