import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pin = el.querySelector("g");
  animate(pin, { translateY: [-6, 1, 0], scaleY: [1.15, 0.9, 1] }, { duration: τmid });

  const plus = el.querySelectorAll('path[d="M16 18h6"], path[d="M19 15v6"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
