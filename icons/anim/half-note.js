import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const note = el.querySelector("g");
  animate(note, { rotate: [0, -12, 8, 0] }, { duration: τmax });

  const head = el.querySelector("circle");
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
