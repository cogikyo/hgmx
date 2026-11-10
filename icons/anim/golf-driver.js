import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const club = el.querySelector("g");
  animate(club, { rotate: [0, 6, -3, 0] }, { duration: τmid });
  const ball = el.querySelector("circle");
  animate(ball, { translateY: [0, -2.5, 0.5, 0] }, { duration: τmid, delay: τmin });
}
