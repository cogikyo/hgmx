import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pin = el.querySelector("g");
  animate(pin, { translateY: [-6, 1, 0], scaleY: [1.15, 0.9, 1] }, { duration: τmid });

  const lid = el.querySelector("rect");
  animate(lid, { scaleY: [0.25, 1.1, 1], opacity: [0.5, 1] }, { duration: τmid, delay: τmin });
}
