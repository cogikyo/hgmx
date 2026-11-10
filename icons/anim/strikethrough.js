import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const letter = el.querySelectorAll("path");
  animate(letter, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const strike = el.querySelector("line");
  animate(strike, { scaleX: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
