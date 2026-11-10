import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("g");
  const van = [body, ...el.querySelectorAll("circle")];
  animate(van, { translateX: [-6, 1, 0], opacity: [0, 1] }, { duration: τmax });
  animate(body, { rotate: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
