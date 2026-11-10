import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const face = el.querySelector("rect");
  animate(face, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const note = el.querySelector("g");
  animate(note, { rotate: [0, -12, 8, 0] }, { duration: τmid, delay: τmin });
}
