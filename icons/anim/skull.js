import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const skull = el.querySelector("g");
  animate(skull, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
  const eyes = el.querySelectorAll("circle");
  animate(eyes, { scale: [1, 1.5, 1] }, { duration: τmin, delay: τmid });
}
