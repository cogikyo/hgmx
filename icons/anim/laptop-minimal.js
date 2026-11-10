import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lid = el.querySelector("rect");
  animate(lid, { scaleY: [0.25, 1.1, 1], opacity: [0.5, 1] }, { duration: τmid });

  const base = el.querySelector("line");
  animate(base, { scaleX: [1, 1.15, 1] }, { duration: τmin, delay: τmid });
}
