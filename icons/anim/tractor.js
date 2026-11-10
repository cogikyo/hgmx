import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("g");
  const tractor = [body, ...el.querySelectorAll(":scope > path, circle")];
  animate(tractor, { translateX: [-6, 1, 0], opacity: [0, 1] }, { duration: τmax });

  animate(body, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
